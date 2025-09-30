package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DBStreamer handles database change notifications
type DBStreamer struct {
	db       *gorm.DB
	hub      *Hub
	listener *pq.Listener
	logger   *zap.Logger
	tables   map[string]bool // Tables to watch
}

// DBChangeEvent represents a database change event
type DBChangeEvent struct {
	Table     string                 `json:"table"`
	Operation string                 `json:"operation"` // INSERT, UPDATE, DELETE
	Data      map[string]interface{} `json:"data"`
	OldData   map[string]interface{} `json:"old_data,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// NewDBStreamer creates a new database streamer
func NewDBStreamer(db *gorm.DB, hub *Hub, logger *zap.Logger) *DBStreamer {
	return &DBStreamer{
		db:     db,
		hub:    hub,
		logger: logger,
		tables: make(map[string]bool),
	}
}

// WatchTable adds a table to watch for changes
func (s *DBStreamer) WatchTable(tableName string) error {
	s.tables[tableName] = true

	// Create trigger function if not exists
	triggerFunc := fmt.Sprintf(`
		CREATE OR REPLACE FUNCTION notify_%s_changes()
		RETURNS trigger AS $$
		DECLARE
			data json;
			notification json;
		BEGIN
			IF (TG_OP = 'DELETE') THEN
				data = row_to_json(OLD);
			ELSE
				data = row_to_json(NEW);
			END IF;
			
			notification = json_build_object(
				'table', TG_TABLE_NAME,
				'operation', TG_OP,
				'data', data,
				'timestamp', NOW()
			);
			
			PERFORM pg_notify('db_changes', notification::text);
			RETURN NULL;
		END;
		$$ LANGUAGE plpgsql;
	`, tableName)

	if err := s.db.Exec(triggerFunc).Error; err != nil {
		return fmt.Errorf("failed to create trigger function: %w", err)
	}

	// Create triggers for INSERT, UPDATE, DELETE
	for _, op := range []string{"INSERT", "UPDATE", "DELETE"} {
		triggerName := fmt.Sprintf("%s_%s_notify", tableName, op)
		trigger := fmt.Sprintf(`
			DROP TRIGGER IF EXISTS %s ON %s;
			CREATE TRIGGER %s
			AFTER %s ON %s
			FOR EACH ROW EXECUTE FUNCTION notify_%s_changes();
		`, triggerName, tableName, triggerName, op, tableName, tableName)

		if err := s.db.Exec(trigger).Error; err != nil {
			return fmt.Errorf("failed to create %s trigger: %w", op, err)
		}
	}

	s.logger.Info("Watching table for changes", zap.String("table", tableName))
	return nil
}

// Start begins listening for database changes
func (s *DBStreamer) Start(ctx context.Context) error {
	// For PostgreSQL, we need to construct connection string
	// This is a simplified version - in production, get from config
	connStr := "host=localhost port=5433 user=app password=secret dbname=myapp sslmode=disable"

	// Create listener
	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			s.logger.Error("Listener error", zap.Error(err))
		}
	}

	s.listener = pq.NewListener(connStr, 10*time.Second, time.Minute, reportProblem)

	// Listen on channel
	if err := s.listener.Listen("db_changes"); err != nil {
		return fmt.Errorf("failed to listen on channel: %w", err)
	}

	s.logger.Info("Database change listener started")

	// Start listening for notifications
	go s.listenLoop(ctx)

	return nil
}

func (s *DBStreamer) listenLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Stopping database change listener")
			s.listener.Close()
			return

		case notification := <-s.listener.Notify:
			if notification == nil {
				continue
			}

			// Parse notification
			var event DBChangeEvent
			if err := json.Unmarshal([]byte(notification.Extra), &event); err != nil {
				s.logger.Error("Failed to parse notification", zap.Error(err))
				continue
			}

			// Broadcast to WebSocket clients
			s.broadcastDBChange(&event)

		case <-time.After(90 * time.Second):
			// Ping to check connection
			if err := s.listener.Ping(); err != nil {
				s.logger.Error("Listener ping failed", zap.Error(err))
			}
		}
	}
}

func (s *DBStreamer) broadcastDBChange(event *DBChangeEvent) {
	// Only broadcast if table is being watched
	if !s.tables[event.Table] {
		return
	}

	s.logger.Debug("Broadcasting database change",
		zap.String("table", event.Table),
		zap.String("operation", event.Operation),
	)

	// Broadcast to specific table channel
	channelName := fmt.Sprintf("db:%s", event.Table)
	payload := map[string]interface{}{
		"table":     event.Table,
		"operation": event.Operation,
		"data":      event.Data,
		"timestamp": event.Timestamp,
	}

	if event.OldData != nil {
		payload["old_data"] = event.OldData
	}

	s.hub.BroadcastToChannel(channelName, "db_change", payload)

	// Also broadcast to general db:* channel
	s.hub.BroadcastToChannel("db:*", "db_change", payload)
}

// Stop stops the database change listener
func (s *DBStreamer) Stop() {
	if s.listener != nil {
		s.listener.Close()
	}
}

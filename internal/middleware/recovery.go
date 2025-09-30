package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery middleware for panic recovery
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if ne, ok := recovered.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
					strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					c.Abort()
					return
				}
			}
		}

		httpRequest, _ := httputil.DumpRequest(c.Request, false)
		if httpRequest != nil {
			logger.Error("Panic recovered",
				zap.Any("error", recovered),
				zap.String("request", string(httpRequest)),
				zap.String("stack", string(debug.Stack())),
			)
		} else {
			logger.Error("Panic recovered",
				zap.Any("error", recovered),
				zap.String("stack", string(debug.Stack())),
			)
		}

		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

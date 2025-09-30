package files

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/db/repository"
	"go-mobile-backend-template/internal/services/storage"
	"go-mobile-backend-template/internal/utils"
	"go-mobile-backend-template/pkg/config"
)

// Handler handles file requests
type Handler struct {
	fileRepo repository.FileRepository
	r2Client *storage.R2Client
	logger   *zap.Logger
	cfg      *config.Config
}

// NewHandler creates a new file handler
func NewHandler(db *gorm.DB, logger *zap.Logger, cfg *config.Config) *Handler {
	r2Client, err := storage.NewR2Client(storage.R2Config{
		AccountID: cfg.R2.AccountID,
		AccessKey: cfg.R2.AccessKey,
		SecretKey: cfg.R2.SecretKey,
		Bucket:    cfg.R2.Bucket,
		Endpoint:  cfg.R2.Endpoint,
	})
	if err != nil {
		logger.Fatal("Failed to create R2 client", zap.Error(err))
	}

	return &Handler{
		fileRepo: repository.NewFileRepository(db),
		r2Client: r2Client,
		logger:   logger,
		cfg:      cfg,
	}
}

// Upload handles file upload
// @Summary Upload file
// @Description Upload a file to R2 storage
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "File to upload"
// @Success 201 {object} UploadResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /files/upload [post]
func (h *Handler) Upload(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "File is required")
		return
	}
	defer file.Close()

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		h.logger.Error("Failed to read file", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to read file")
		return
	}

	// Generate unique key for R2
	ext := filepath.Ext(header.Filename)
	r2Key := fmt.Sprintf("uploads/%d/%s%s", userID.(uint), uuid.New().String(), ext)

	// Upload to R2
	ctx := context.Background()
	r2URL, err := h.r2Client.Upload(ctx, r2Key, fileContent, header.Header.Get("Content-Type"))
	if err != nil {
		h.logger.Error("Failed to upload file to R2", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to upload file")
		return
	}

	// Save file metadata
	fileRecord := &repository.File{
		UserID:   userID.(uint),
		FileName: header.Filename,
		FileSize: header.Size,
		FileType: header.Header.Get("Content-Type"),
		R2Key:    r2Key,
		R2URL:    r2URL,
		IsPublic: false,
	}

	if err := h.fileRepo.Create(ctx, fileRecord); err != nil {
		h.logger.Error("Failed to save file metadata", zap.Error(err))
		// Try to delete from R2
		h.r2Client.Delete(ctx, r2Key)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save file")
		return
	}

	response := UploadResponse{
		File: FileResponse{
			ID:        fileRecord.ID,
			FileName:  fileRecord.FileName,
			FileSize:  fileRecord.FileSize,
			FileType:  fileRecord.FileType,
			R2URL:     fileRecord.R2URL,
			IsPublic:  fileRecord.IsPublic,
			CreatedAt: fileRecord.CreatedAt.Format(time.RFC3339),
		},
	}

	utils.SuccessResponse(c, http.StatusCreated, "File uploaded successfully", response)
}

// GetFile gets file metadata
// @Summary Get file metadata
// @Description Get file metadata by ID
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "File ID"
// @Success 200 {object} FileResponse
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /files/{id} [get]
func (h *Handler) GetFile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid file ID")
		return
	}

	ctx := context.Background()
	file, err := h.fileRepo.GetByID(ctx, uint(fileID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "File not found")
		return
	}

	// Check ownership
	if file.UserID != userID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied")
		return
	}

	response := FileResponse{
		ID:        file.ID,
		FileName:  file.FileName,
		FileSize:  file.FileSize,
		FileType:  file.FileType,
		R2URL:     file.R2URL,
		IsPublic:  file.IsPublic,
		CreatedAt: file.CreatedAt.Format(time.RFC3339),
	}

	utils.SuccessResponse(c, http.StatusOK, "File retrieved successfully", response)
}

// GetDownloadURL generates a presigned download URL
// @Summary Get download URL
// @Description Generate a presigned URL for downloading a file
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "File ID"
// @Success 200 {object} DownloadURLResponse
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /files/{id}/download [get]
func (h *Handler) GetDownloadURL(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid file ID")
		return
	}

	ctx := context.Background()
	file, err := h.fileRepo.GetByID(ctx, uint(fileID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "File not found")
		return
	}

	// Check ownership
	if file.UserID != userID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied")
		return
	}

	// Generate presigned URL (valid for 1 hour)
	expiration := 1 * time.Hour
	downloadURL, err := h.r2Client.GeneratePresignedURL(ctx, file.R2Key, expiration)
	if err != nil {
		h.logger.Error("Failed to generate download URL", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate download URL")
		return
	}

	response := DownloadURLResponse{
		URL:       downloadURL,
		ExpiresIn: int(expiration.Seconds()),
	}

	utils.SuccessResponse(c, http.StatusOK, "Download URL generated successfully", response)
}

// DeleteFile deletes a file
// @Summary Delete file
// @Description Delete a file by ID
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "File ID"
// @Success 200 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /files/{id} [delete]
func (h *Handler) DeleteFile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	fileID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid file ID")
		return
	}

	ctx := context.Background()
	file, err := h.fileRepo.GetByID(ctx, uint(fileID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "File not found")
		return
	}

	// Check ownership
	if file.UserID != userID.(uint) {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied")
		return
	}

	// Delete from R2
	if err := h.r2Client.Delete(ctx, file.R2Key); err != nil {
		h.logger.Error("Failed to delete file from R2", zap.Error(err))
		// Continue with database deletion even if R2 deletion fails
	}

	// Delete from database
	if err := h.fileRepo.Delete(ctx, uint(fileID)); err != nil {
		h.logger.Error("Failed to delete file from database", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete file")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "File deleted successfully", nil)
}

// ListFiles lists user's files
// @Summary List files
// @Description List current user's files
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} FileResponse
// @Failure 401 {object} utils.Response
// @Router /files [get]
func (h *Handler) ListFiles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	ctx := context.Background()
	files, err := h.fileRepo.GetByUserID(ctx, userID.(uint), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list files", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to list files")
		return
	}

	var response []FileResponse
	for _, file := range files {
		response = append(response, FileResponse{
			ID:        file.ID,
			FileName:  file.FileName,
			FileSize:  file.FileSize,
			FileType:  file.FileType,
			R2URL:     file.R2URL,
			IsPublic:  file.IsPublic,
			CreatedAt: file.CreatedAt.Format(time.RFC3339),
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Files retrieved successfully", response)
}

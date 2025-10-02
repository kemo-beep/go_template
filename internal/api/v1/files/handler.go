package files

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-mobile-backend-template/internal/db/repository/generated"
	"go-mobile-backend-template/internal/utils"
)

// Handler handles files requests
type Handler struct {
	filesRepo generated.FilesRepository
	logger             *zap.Logger
}

// NewHandler creates a new files handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		filesRepo: generated.NewFilesRepository(db),
		logger:             logger,
	}
}

// CreateFiles creates a new files
// @Summary Create files
// @Description Create a new files record
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body FilesCreateRequest true "Create files request"
// @Success 201 {object} FilesResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /files [post]
func (h *Handler) CreateFiles(c *gin.Context) {
	var req FilesCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	files := &generated.Files{

		Userid: req.Userid,

		Filename: req.Filename,

		Filesize: req.Filesize,

		Filetype: req.Filetype,

		R2key: req.R2key,

		R2url: req.R2url,

		Ispublic: req.Ispublic,

		Deletedat: req.Deletedat,

	}

	ctx := context.Background()
	if err := h.filesRepo.Create(ctx, files); err != nil {
		h.logger.Error("Failed to create files", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create files")
		return
	}

	response := FilesResponse{

		Id: files.Id,

		Userid: files.Userid,

		Filename: files.Filename,

		Filesize: files.Filesize,

		Filetype: files.Filetype,

		R2key: files.R2key,

		R2url: files.R2url,

		Ispublic: files.Ispublic,

		Createdat: files.Createdat,

		Updatedat: files.Updatedat,

		Deletedat: files.Deletedat,

	}

	utils.SuccessResponse(c, http.StatusCreated, "files created successfully", response)
}

// GetFiles gets a files by ID
// @Summary Get files
// @Description Get a files by ID
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "files ID"
// @Success 200 {object} FilesResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /files/{id} [get]
func (h *Handler) GetFiles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	files, err := h.filesRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "files not found")
			return
		}
		h.logger.Error("Failed to get files", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get files")
		return
	}

	response := FilesResponse{

		Id: files.Id,

		Userid: files.Userid,

		Filename: files.Filename,

		Filesize: files.Filesize,

		Filetype: files.Filetype,

		R2key: files.R2key,

		R2url: files.R2url,

		Ispublic: files.Ispublic,

		Createdat: files.Createdat,

		Updatedat: files.Updatedat,

		Deletedat: files.Deletedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "files retrieved successfully", response)
}

// GetAllFiless gets all filess with pagination
// @Summary Get all filess
// @Description Get all filess with pagination
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /files [get]
func (h *Handler) GetAllFiless(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	ctx := context.Background()
	filess, total, err := h.filesRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get filess", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get filess")
		return
	}

	var responses []FilesResponse
	for _, files := range filess {
		responses = append(responses, FilesResponse{

			Id: files.Id,

			Userid: files.Userid,

			Filename: files.Filename,

			Filesize: files.Filesize,

			Filetype: files.Filetype,

			R2key: files.R2key,

			R2url: files.R2url,

			Ispublic: files.Ispublic,

			Createdat: files.Createdat,

			Updatedat: files.Updatedat,

			Deletedat: files.Deletedat,

		})
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	pagination := PaginationResponse{
		Data: responses,
		Pagination: PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "filess retrieved successfully", pagination)
}

// UpdateFiles updates a files
// @Summary Update files
// @Description Update a files by ID
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "files ID"
// @Param request body FilesUpdateRequest true "Update files request"
// @Success 200 {object} FilesResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /files/{id} [put]
func (h *Handler) UpdateFiles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req FilesUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	files, err := h.filesRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "files not found")
			return
		}
		h.logger.Error("Failed to get files", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get files")
		return
	}


	
	files.Userid = req.Userid
	

	
	files.Filename = req.Filename
	

	
	files.Filesize = req.Filesize
	

	
	files.Filetype = req.Filetype
	

	
	files.R2key = req.R2key
	

	
	files.R2url = req.R2url
	

	
	files.Ispublic = req.Ispublic
	

	
	files.Deletedat = req.Deletedat
	


	if err := h.filesRepo.Update(ctx, files); err != nil {
		h.logger.Error("Failed to update files", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update files")
		return
	}

	response := FilesResponse{

		Id: files.Id,

		Userid: files.Userid,

		Filename: files.Filename,

		Filesize: files.Filesize,

		Filetype: files.Filetype,

		R2key: files.R2key,

		R2url: files.R2url,

		Ispublic: files.Ispublic,

		Createdat: files.Createdat,

		Updatedat: files.Updatedat,

		Deletedat: files.Deletedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "files updated successfully", response)
}

// DeleteFiles deletes a files
// @Summary Delete files
// @Description Delete a files by ID
// @Tags files
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "files ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /files/{id} [delete]
func (h *Handler) DeleteFiles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.filesRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "files not found")
			return
		}
		h.logger.Error("Failed to delete files", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete files")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "files deleted successfully", nil)
}

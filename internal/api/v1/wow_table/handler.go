package wow_table

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

// Handler handles wow_table requests
type Handler struct {
	wowTableRepo generated.WowtableRepository
	logger             *zap.Logger
}

// NewHandler creates a new wow_table handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		wowTableRepo: generated.NewWowtableRepository(db),
		logger:             logger,
	}
}

// CreateWowtable creates a new wow_table
// @Summary Create wow_table
// @Description Create a new wow_table record
// @Tags wow_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body WowtableCreateRequest true "Create wow_table request"
// @Success 201 {object} WowtableResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /wow_table [post]
func (h *Handler) CreateWowtable(c *gin.Context) {
	var req WowtableCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	wowTable := &generated.Wowtable{

		Swim: req.Swim,

	}

	ctx := context.Background()
	if err := h.wowTableRepo.Create(ctx, wowTable); err != nil {
		h.logger.Error("Failed to create wow_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create wow_table")
		return
	}

	response := WowtableResponse{

		Id: wowTable.Id,

		Swim: wowTable.Swim,

	}

	utils.SuccessResponse(c, http.StatusCreated, "wow_table created successfully", response)
}

// GetWowtable gets a wow_table by ID
// @Summary Get wow_table
// @Description Get a wow_table by ID
// @Tags wow_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "wow_table ID"
// @Success 200 {object} WowtableResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /wow_table/{id} [get]
func (h *Handler) GetWowtable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	wowTable, err := h.wowTableRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "wow_table not found")
			return
		}
		h.logger.Error("Failed to get wow_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get wow_table")
		return
	}

	response := WowtableResponse{

		Id: wowTable.Id,

		Swim: wowTable.Swim,

	}

	utils.SuccessResponse(c, http.StatusOK, "wow_table retrieved successfully", response)
}

// GetAllWowtables gets all wow_tables with pagination
// @Summary Get all wow_tables
// @Description Get all wow_tables with pagination
// @Tags wow_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /wow_table [get]
func (h *Handler) GetAllWowtables(c *gin.Context) {
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
	wowTables, total, err := h.wowTableRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get wow_tables", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get wow_tables")
		return
	}

	var responses []WowtableResponse
	for _, wowTable := range wowTables {
		responses = append(responses, WowtableResponse{

			Id: wowTable.Id,

			Swim: wowTable.Swim,

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

	utils.SuccessResponse(c, http.StatusOK, "wow_tables retrieved successfully", pagination)
}

// UpdateWowtable updates a wow_table
// @Summary Update wow_table
// @Description Update a wow_table by ID
// @Tags wow_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "wow_table ID"
// @Param request body WowtableUpdateRequest true "Update wow_table request"
// @Success 200 {object} WowtableResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /wow_table/{id} [put]
func (h *Handler) UpdateWowtable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req WowtableUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	wowTable, err := h.wowTableRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "wow_table not found")
			return
		}
		h.logger.Error("Failed to get wow_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get wow_table")
		return
	}


	
	wowTable.Swim = req.Swim
	


	if err := h.wowTableRepo.Update(ctx, wowTable); err != nil {
		h.logger.Error("Failed to update wow_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update wow_table")
		return
	}

	response := WowtableResponse{

		Id: wowTable.Id,

		Swim: wowTable.Swim,

	}

	utils.SuccessResponse(c, http.StatusOK, "wow_table updated successfully", response)
}

// DeleteWowtable deletes a wow_table
// @Summary Delete wow_table
// @Description Delete a wow_table by ID
// @Tags wow_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "wow_table ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /wow_table/{id} [delete]
func (h *Handler) DeleteWowtable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.wowTableRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "wow_table not found")
			return
		}
		h.logger.Error("Failed to delete wow_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete wow_table")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "wow_table deleted successfully", nil)
}

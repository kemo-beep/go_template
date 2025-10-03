package dancing_table

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

// Handler handles dancing_table requests
type Handler struct {
	dancingTableRepo generated.DancingtableRepository
	logger             *zap.Logger
}

// NewHandler creates a new dancing_table handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		dancingTableRepo: generated.NewDancingtableRepository(db),
		logger:             logger,
	}
}

// CreateDancingtable creates a new dancing_table
// @Summary Create dancing_table
// @Description Create a new dancing_table record
// @Tags dancing_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body DancingtableCreateRequest true "Create dancing_table request"
// @Success 201 {object} DancingtableResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /dancing_table [post]
func (h *Handler) CreateDancingtable(c *gin.Context) {
	var req DancingtableCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	dancingTable := &generated.Dancingtable{

		Frequency: req.Frequency,

	}

	ctx := context.Background()
	if err := h.dancingTableRepo.Create(ctx, dancingTable); err != nil {
		h.logger.Error("Failed to create dancing_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create dancing_table")
		return
	}

	response := DancingtableResponse{

		Id: dancingTable.Id,

		Frequency: dancingTable.Frequency,

	}

	utils.SuccessResponse(c, http.StatusCreated, "dancing_table created successfully", response)
}

// GetDancingtable gets a dancing_table by ID
// @Summary Get dancing_table
// @Description Get a dancing_table by ID
// @Tags dancing_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "dancing_table ID"
// @Success 200 {object} DancingtableResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /dancing_table/{id} [get]
func (h *Handler) GetDancingtable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	dancingTable, err := h.dancingTableRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "dancing_table not found")
			return
		}
		h.logger.Error("Failed to get dancing_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get dancing_table")
		return
	}

	response := DancingtableResponse{

		Id: dancingTable.Id,

		Frequency: dancingTable.Frequency,

	}

	utils.SuccessResponse(c, http.StatusOK, "dancing_table retrieved successfully", response)
}

// GetAllDancingtables gets all dancing_tables with pagination
// @Summary Get all dancing_tables
// @Description Get all dancing_tables with pagination
// @Tags dancing_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /dancing_table [get]
func (h *Handler) GetAllDancingtables(c *gin.Context) {
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
	dancingTables, total, err := h.dancingTableRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get dancing_tables", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get dancing_tables")
		return
	}

	var responses []DancingtableResponse
	for _, dancingTable := range dancingTables {
		responses = append(responses, DancingtableResponse{

			Id: dancingTable.Id,

			Frequency: dancingTable.Frequency,

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

	utils.SuccessResponse(c, http.StatusOK, "dancing_tables retrieved successfully", pagination)
}

// UpdateDancingtable updates a dancing_table
// @Summary Update dancing_table
// @Description Update a dancing_table by ID
// @Tags dancing_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "dancing_table ID"
// @Param request body DancingtableUpdateRequest true "Update dancing_table request"
// @Success 200 {object} DancingtableResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /dancing_table/{id} [put]
func (h *Handler) UpdateDancingtable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req DancingtableUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	dancingTable, err := h.dancingTableRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "dancing_table not found")
			return
		}
		h.logger.Error("Failed to get dancing_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get dancing_table")
		return
	}


	
	dancingTable.Frequency = req.Frequency
	


	if err := h.dancingTableRepo.Update(ctx, dancingTable); err != nil {
		h.logger.Error("Failed to update dancing_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update dancing_table")
		return
	}

	response := DancingtableResponse{

		Id: dancingTable.Id,

		Frequency: dancingTable.Frequency,

	}

	utils.SuccessResponse(c, http.StatusOK, "dancing_table updated successfully", response)
}

// DeleteDancingtable deletes a dancing_table
// @Summary Delete dancing_table
// @Description Delete a dancing_table by ID
// @Tags dancing_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "dancing_table ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /dancing_table/{id} [delete]
func (h *Handler) DeleteDancingtable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.dancingTableRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "dancing_table not found")
			return
		}
		h.logger.Error("Failed to delete dancing_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete dancing_table")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "dancing_table deleted successfully", nil)
}

package test_table

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

// Handler handles test_table requests
type Handler struct {
	testTableRepo generated.TesttableRepository
	logger             *zap.Logger
}

// NewHandler creates a new test_table handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		testTableRepo: generated.NewTesttableRepository(db),
		logger:             logger,
	}
}

// CreateTesttable creates a new test_table
// @Summary Create test_table
// @Description Create a new test_table record
// @Tags test_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TesttableCreateRequest true "Create test_table request"
// @Success 201 {object} TesttableResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /test_table [post]
func (h *Handler) CreateTesttable(c *gin.Context) {
	var req TesttableCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	testTable := &generated.Testtable{

		Name: req.Name,

		Familyname: req.Familyname,

		Prefrence: req.Prefrence,

		Preferences: req.Preferences,

	}

	ctx := context.Background()
	if err := h.testTableRepo.Create(ctx, testTable); err != nil {
		h.logger.Error("Failed to create test_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create test_table")
		return
	}

	response := TesttableResponse{

		Id: testTable.Id,

		Name: testTable.Name,

		Familyname: testTable.Familyname,

		Prefrence: testTable.Prefrence,

		Preferences: testTable.Preferences,

	}

	utils.SuccessResponse(c, http.StatusCreated, "test_table created successfully", response)
}

// GetTesttable gets a test_table by ID
// @Summary Get test_table
// @Description Get a test_table by ID
// @Tags test_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "test_table ID"
// @Success 200 {object} TesttableResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /test_table/{id} [get]
func (h *Handler) GetTesttable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	testTable, err := h.testTableRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "test_table not found")
			return
		}
		h.logger.Error("Failed to get test_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get test_table")
		return
	}

	response := TesttableResponse{

		Id: testTable.Id,

		Name: testTable.Name,

		Familyname: testTable.Familyname,

		Prefrence: testTable.Prefrence,

		Preferences: testTable.Preferences,

	}

	utils.SuccessResponse(c, http.StatusOK, "test_table retrieved successfully", response)
}

// GetAllTesttables gets all test_tables with pagination
// @Summary Get all test_tables
// @Description Get all test_tables with pagination
// @Tags test_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /test_table [get]
func (h *Handler) GetAllTesttables(c *gin.Context) {
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
	testTables, total, err := h.testTableRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get test_tables", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get test_tables")
		return
	}

	var responses []TesttableResponse
	for _, testTable := range testTables {
		responses = append(responses, TesttableResponse{

			Id: testTable.Id,

			Name: testTable.Name,

			Familyname: testTable.Familyname,

			Prefrence: testTable.Prefrence,

			Preferences: testTable.Preferences,

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

	utils.SuccessResponse(c, http.StatusOK, "test_tables retrieved successfully", pagination)
}

// UpdateTesttable updates a test_table
// @Summary Update test_table
// @Description Update a test_table by ID
// @Tags test_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "test_table ID"
// @Param request body TesttableUpdateRequest true "Update test_table request"
// @Success 200 {object} TesttableResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /test_table/{id} [put]
func (h *Handler) UpdateTesttable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req TesttableUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	testTable, err := h.testTableRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "test_table not found")
			return
		}
		h.logger.Error("Failed to get test_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get test_table")
		return
	}


	
	testTable.Name = req.Name
	

	
	testTable.Familyname = req.Familyname
	

	
	testTable.Prefrence = req.Prefrence
	

	
	testTable.Preferences = req.Preferences
	


	if err := h.testTableRepo.Update(ctx, testTable); err != nil {
		h.logger.Error("Failed to update test_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update test_table")
		return
	}

	response := TesttableResponse{

		Id: testTable.Id,

		Name: testTable.Name,

		Familyname: testTable.Familyname,

		Prefrence: testTable.Prefrence,

		Preferences: testTable.Preferences,

	}

	utils.SuccessResponse(c, http.StatusOK, "test_table updated successfully", response)
}

// DeleteTesttable deletes a test_table
// @Summary Delete test_table
// @Description Delete a test_table by ID
// @Tags test_table
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "test_table ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /test_table/{id} [delete]
func (h *Handler) DeleteTesttable(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.testTableRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "test_table not found")
			return
		}
		h.logger.Error("Failed to delete test_table", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete test_table")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "test_table deleted successfully", nil)
}

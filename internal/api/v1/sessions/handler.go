package sessions

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

// Handler handles sessions requests
type Handler struct {
	sessionsRepo generated.SessionsRepository
	logger             *zap.Logger
}

// NewHandler creates a new sessions handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		sessionsRepo: generated.NewSessionsRepository(db),
		logger:             logger,
	}
}

// CreateSessions creates a new sessions
// @Summary Create sessions
// @Description Create a new sessions record
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SessionsCreateRequest true "Create sessions request"
// @Success 201 {object} SessionsResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /sessions [post]
func (h *Handler) CreateSessions(c *gin.Context) {
	var req SessionsCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	sessions := &generated.Sessions{

		Userid: req.Userid,

		Token: req.Token,

		Refreshtoken: req.Refreshtoken,

		Deviceinfo: req.Deviceinfo,

		Ipaddress: req.Ipaddress,

		Useragent: req.Useragent,

		Isactive: req.Isactive,

		Expiresat: req.Expiresat,

		Lastusedat: req.Lastusedat,

	}

	ctx := context.Background()
	if err := h.sessionsRepo.Create(ctx, sessions); err != nil {
		h.logger.Error("Failed to create sessions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create sessions")
		return
	}

	response := SessionsResponse{

		Id: sessions.Id,

		Userid: sessions.Userid,

		Token: sessions.Token,

		Refreshtoken: sessions.Refreshtoken,

		Deviceinfo: sessions.Deviceinfo,

		Ipaddress: sessions.Ipaddress,

		Useragent: sessions.Useragent,

		Isactive: sessions.Isactive,

		Expiresat: sessions.Expiresat,

		Createdat: sessions.Createdat,

		Lastusedat: sessions.Lastusedat,

	}

	utils.SuccessResponse(c, http.StatusCreated, "sessions created successfully", response)
}

// GetSessions gets a sessions by ID
// @Summary Get sessions
// @Description Get a sessions by ID
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "sessions ID"
// @Success 200 {object} SessionsResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /sessions/{id} [get]
func (h *Handler) GetSessions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	sessions, err := h.sessionsRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "sessions not found")
			return
		}
		h.logger.Error("Failed to get sessions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get sessions")
		return
	}

	response := SessionsResponse{

		Id: sessions.Id,

		Userid: sessions.Userid,

		Token: sessions.Token,

		Refreshtoken: sessions.Refreshtoken,

		Deviceinfo: sessions.Deviceinfo,

		Ipaddress: sessions.Ipaddress,

		Useragent: sessions.Useragent,

		Isactive: sessions.Isactive,

		Expiresat: sessions.Expiresat,

		Createdat: sessions.Createdat,

		Lastusedat: sessions.Lastusedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "sessions retrieved successfully", response)
}

// GetAllSessionss gets all sessionss with pagination
// @Summary Get all sessionss
// @Description Get all sessionss with pagination
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /sessions [get]
func (h *Handler) GetAllSessionss(c *gin.Context) {
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
	sessionss, total, err := h.sessionsRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get sessionss", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get sessionss")
		return
	}

	var responses []SessionsResponse
	for _, sessions := range sessionss {
		responses = append(responses, SessionsResponse{

			Id: sessions.Id,

			Userid: sessions.Userid,

			Token: sessions.Token,

			Refreshtoken: sessions.Refreshtoken,

			Deviceinfo: sessions.Deviceinfo,

			Ipaddress: sessions.Ipaddress,

			Useragent: sessions.Useragent,

			Isactive: sessions.Isactive,

			Expiresat: sessions.Expiresat,

			Createdat: sessions.Createdat,

			Lastusedat: sessions.Lastusedat,

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

	utils.SuccessResponse(c, http.StatusOK, "sessionss retrieved successfully", pagination)
}

// UpdateSessions updates a sessions
// @Summary Update sessions
// @Description Update a sessions by ID
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "sessions ID"
// @Param request body SessionsUpdateRequest true "Update sessions request"
// @Success 200 {object} SessionsResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /sessions/{id} [put]
func (h *Handler) UpdateSessions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req SessionsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	sessions, err := h.sessionsRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "sessions not found")
			return
		}
		h.logger.Error("Failed to get sessions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get sessions")
		return
	}


	
	sessions.Userid = req.Userid
	

	
	sessions.Token = req.Token
	

	
	sessions.Refreshtoken = req.Refreshtoken
	

	
	sessions.Deviceinfo = req.Deviceinfo
	

	
	sessions.Ipaddress = req.Ipaddress
	

	
	sessions.Useragent = req.Useragent
	

	
	sessions.Isactive = req.Isactive
	

	
	sessions.Expiresat = req.Expiresat
	

	
	sessions.Lastusedat = req.Lastusedat
	


	if err := h.sessionsRepo.Update(ctx, sessions); err != nil {
		h.logger.Error("Failed to update sessions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update sessions")
		return
	}

	response := SessionsResponse{

		Id: sessions.Id,

		Userid: sessions.Userid,

		Token: sessions.Token,

		Refreshtoken: sessions.Refreshtoken,

		Deviceinfo: sessions.Deviceinfo,

		Ipaddress: sessions.Ipaddress,

		Useragent: sessions.Useragent,

		Isactive: sessions.Isactive,

		Expiresat: sessions.Expiresat,

		Createdat: sessions.Createdat,

		Lastusedat: sessions.Lastusedat,

	}

	utils.SuccessResponse(c, http.StatusOK, "sessions updated successfully", response)
}

// DeleteSessions deletes a sessions
// @Summary Delete sessions
// @Description Delete a sessions by ID
// @Tags sessions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "sessions ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /sessions/{id} [delete]
func (h *Handler) DeleteSessions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.sessionsRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "sessions not found")
			return
		}
		h.logger.Error("Failed to delete sessions", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete sessions")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "sessions deleted successfully", nil)
}

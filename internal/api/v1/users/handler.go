package users

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

// Handler handles users requests
type Handler struct {
	usersRepo generated.UsersRepository
	logger             *zap.Logger
}

// NewHandler creates a new users handler
func NewHandler(db *gorm.DB, logger *zap.Logger) *Handler {
	return &Handler{
		usersRepo: generated.NewUsersRepository(db),
		logger:             logger,
	}
}

// CreateUsers creates a new users
// @Summary Create users
// @Description Create a new users record
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UsersCreateRequest true "Create users request"
// @Success 201 {object} UsersResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /users [post]
func (h *Handler) CreateUsers(c *gin.Context) {
	var req UsersCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	users := &generated.Users{

		Email: req.Email,

		Password: req.Password,

		Name: req.Name,

		Isactive: req.Isactive,

		Isadmin: req.Isadmin,

		Deletedat: req.Deletedat,

		Emailverified: req.Emailverified,

		Emailverifiedat: req.Emailverifiedat,

		Lastloginat: req.Lastloginat,

		Failedloginattempts: req.Failedloginattempts,

		Lockeduntil: req.Lockeduntil,

		Metadata: req.Metadata,

		Nickname: req.Nickname,

		Bio: req.Bio,

	}

	ctx := context.Background()
	if err := h.usersRepo.Create(ctx, users); err != nil {
		h.logger.Error("Failed to create users", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create users")
		return
	}

	response := UsersResponse{

		Id: users.Id,

		Email: users.Email,

		Password: users.Password,

		Name: users.Name,

		Isactive: users.Isactive,

		Isadmin: users.Isadmin,

		Createdat: users.Createdat,

		Updatedat: users.Updatedat,

		Deletedat: users.Deletedat,

		Emailverified: users.Emailverified,

		Emailverifiedat: users.Emailverifiedat,

		Lastloginat: users.Lastloginat,

		Failedloginattempts: users.Failedloginattempts,

		Lockeduntil: users.Lockeduntil,

		Metadata: users.Metadata,

		Nickname: users.Nickname,

		Bio: users.Bio,

	}

	utils.SuccessResponse(c, http.StatusCreated, "users created successfully", response)
}

// GetUsers gets a users by ID
// @Summary Get users
// @Description Get a users by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "users ID"
// @Success 200 {object} UsersResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /users/{id} [get]
func (h *Handler) GetUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	users, err := h.usersRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "users not found")
			return
		}
		h.logger.Error("Failed to get users", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get users")
		return
	}

	response := UsersResponse{

		Id: users.Id,

		Email: users.Email,

		Password: users.Password,

		Name: users.Name,

		Isactive: users.Isactive,

		Isadmin: users.Isadmin,

		Createdat: users.Createdat,

		Updatedat: users.Updatedat,

		Deletedat: users.Deletedat,

		Emailverified: users.Emailverified,

		Emailverifiedat: users.Emailverifiedat,

		Lastloginat: users.Lastloginat,

		Failedloginattempts: users.Failedloginattempts,

		Lockeduntil: users.Lockeduntil,

		Metadata: users.Metadata,

		Nickname: users.Nickname,

		Bio: users.Bio,

	}

	utils.SuccessResponse(c, http.StatusOK, "users retrieved successfully", response)
}

// GetAllUserss gets all userss with pagination
// @Summary Get all userss
// @Description Get all userss with pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(20)
// @Success 200 {object} PaginationResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Router /users [get]
func (h *Handler) GetAllUserss(c *gin.Context) {
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
	userss, total, err := h.usersRepo.GetAll(ctx, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get userss", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get userss")
		return
	}

	var responses []UsersResponse
	for _, users := range userss {
		responses = append(responses, UsersResponse{

			Id: users.Id,

			Email: users.Email,

			Password: users.Password,

			Name: users.Name,

			Isactive: users.Isactive,

			Isadmin: users.Isadmin,

			Createdat: users.Createdat,

			Updatedat: users.Updatedat,

			Deletedat: users.Deletedat,

			Emailverified: users.Emailverified,

			Emailverifiedat: users.Emailverifiedat,

			Lastloginat: users.Lastloginat,

			Failedloginattempts: users.Failedloginattempts,

			Lockeduntil: users.Lockeduntil,

			Metadata: users.Metadata,

			Nickname: users.Nickname,

			Bio: users.Bio,

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

	utils.SuccessResponse(c, http.StatusOK, "userss retrieved successfully", pagination)
}

// UpdateUsers updates a users
// @Summary Update users
// @Description Update a users by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "users ID"
// @Param request body UsersUpdateRequest true "Update users request"
// @Success 200 {object} UsersResponse
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /users/{id} [put]
func (h *Handler) UpdateUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req UsersUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := context.Background()
	users, err := h.usersRepo.GetByID(ctx, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "users not found")
			return
		}
		h.logger.Error("Failed to get users", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get users")
		return
	}


	
	users.Email = req.Email
	

	
	users.Password = req.Password
	

	
	users.Name = req.Name
	

	
	users.Isactive = req.Isactive
	

	
	users.Isadmin = req.Isadmin
	

	
	users.Deletedat = req.Deletedat
	

	
	users.Emailverified = req.Emailverified
	

	
	users.Emailverifiedat = req.Emailverifiedat
	

	
	users.Lastloginat = req.Lastloginat
	

	
	users.Failedloginattempts = req.Failedloginattempts
	

	
	users.Lockeduntil = req.Lockeduntil
	

	
	users.Metadata = req.Metadata
	

	
	users.Nickname = req.Nickname
	

	
	users.Bio = req.Bio
	


	if err := h.usersRepo.Update(ctx, users); err != nil {
		h.logger.Error("Failed to update users", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update users")
		return
	}

	response := UsersResponse{

		Id: users.Id,

		Email: users.Email,

		Password: users.Password,

		Name: users.Name,

		Isactive: users.Isactive,

		Isadmin: users.Isadmin,

		Createdat: users.Createdat,

		Updatedat: users.Updatedat,

		Deletedat: users.Deletedat,

		Emailverified: users.Emailverified,

		Emailverifiedat: users.Emailverifiedat,

		Lastloginat: users.Lastloginat,

		Failedloginattempts: users.Failedloginattempts,

		Lockeduntil: users.Lockeduntil,

		Metadata: users.Metadata,

		Nickname: users.Nickname,

		Bio: users.Bio,

	}

	utils.SuccessResponse(c, http.StatusOK, "users updated successfully", response)
}

// DeleteUsers deletes a users
// @Summary Delete users
// @Description Delete a users by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "users ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 401 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /users/{id} [delete]
func (h *Handler) DeleteUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	ctx := context.Background()
	if err := h.usersRepo.Delete(ctx, uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "users not found")
			return
		}
		h.logger.Error("Failed to delete users", zap.Error(err))
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete users")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "users deleted successfully", nil)
}

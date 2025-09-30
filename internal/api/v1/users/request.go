package users

// UpdateProfileRequest represents update profile request
type UpdateProfileRequest struct {
	Name string `json:"name" binding:"omitempty,min=2"`
}

// ChangePasswordRequest represents change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UserResponse represents user response
type UserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
	IsAdmin  bool   `json:"is_admin"`
}

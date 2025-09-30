package repository

import (
	"context"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
}

// FileRepository defines the interface for file data operations
type FileRepository interface {
	Create(ctx context.Context, file *File) error
	GetByID(ctx context.Context, id uint) (*File, error)
	GetByR2Key(ctx context.Context, r2Key string) (*File, error)
	GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*File, error)
	Update(ctx context.Context, file *File) error
	Delete(ctx context.Context, id uint) error
	DeleteByR2Key(ctx context.Context, r2Key string) error
}

// RefreshTokenRepository defines the interface for refresh token operations
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	GetByToken(ctx context.Context, token string) (*RefreshToken, error)
	GetByUserID(ctx context.Context, userID uint) ([]*RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllForUser(ctx context.Context, userID uint) error
	DeleteExpired(ctx context.Context) error
}

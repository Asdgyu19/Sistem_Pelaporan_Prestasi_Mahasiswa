package app

import (
	"context"
)

// UserRepository defines interface for user data operations
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, id string, user *User) error
}

// AchievementRepository defines interface for achievement data operations
type AchievementRepository interface {
	Create(ctx context.Context, achievement *Achievement) error
	GetByID(ctx context.Context, id string) (*Achievement, error)
	GetByMahasiswaID(ctx context.Context, mahasiswaID string) ([]*Achievement, error)
	GetByStatus(ctx context.Context, status AchievementStatus) ([]*Achievement, error)
	Update(ctx context.Context, id string, achievement *Achievement) error
	SoftDelete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*Achievement, error)
}

// AchievementFileRepository defines interface for file operations (MongoDB)
type AchievementFileRepository interface {
	Create(ctx context.Context, file *AchievementFile) error
	GetByAchievementID(ctx context.Context, achievementID string) ([]*AchievementFile, error)
	Delete(ctx context.Context, id string) error
}

// AuthUsecase defines interface for authentication business logic
type AuthUsecase interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
}

// AchievementUsecase defines interface for achievement business logic
type AchievementUsecase interface {
	Create(ctx context.Context, req *CreateAchievementRequest, mahasiswaID string) (*Achievement, error)
	GetByMahasiswa(ctx context.Context, mahasiswaID string) ([]*Achievement, error)
	GetByID(ctx context.Context, id string) (*Achievement, error)
	Update(ctx context.Context, id string, req *UpdateAchievementRequest, userID string) (*Achievement, error)
	Delete(ctx context.Context, id string, userID string) error
	Submit(ctx context.Context, id string, userID string) error
	Verify(ctx context.Context, id string, verifierID string) error
	Reject(ctx context.Context, id string, verifierID string, reason string) error
	GetPendingVerifications(ctx context.Context) ([]*Achievement, error)
	GetAll(ctx context.Context) ([]*Achievement, error)
}

// FileUsecase defines interface for file operations business logic
type FileUsecase interface {
	UploadAchievementFile(ctx context.Context, achievementID string, fileName string, fileData []byte, mimeType string) (*AchievementFile, error)
	GetAchievementFiles(ctx context.Context, achievementID string) ([]*AchievementFile, error)
	DeleteFile(ctx context.Context, fileID string) error
}

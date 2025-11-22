package app

import (
	"time"
)

// UserRole represents user roles in the system
type UserRole string

const (
	RoleMahasiswa UserRole = "mahasiswa"
	RoleDosenWali UserRole = "dosen_wali"
	RoleAdmin     UserRole = "admin"
)

// AchievementStatus represents achievement status
type AchievementStatus string

const (
	StatusDraft     AchievementStatus = "draft"
	StatusSubmitted AchievementStatus = "submitted"
	StatusVerified  AchievementStatus = "verified"
	StatusRejected  AchievementStatus = "rejected"
)

// User represents user entity
type User struct {
	ID        string    `json:"id" db:"id"`
	NIM       *string   `json:"nim,omitempty" db:"nim"` // untuk mahasiswa
	NIP       *string   `json:"nip,omitempty" db:"nip"` // untuk dosen & admin
	Name      string    `json:"name" db:"name" validate:"required"`
	Email     string    `json:"email" db:"email" validate:"required,email"`
	Password  string    `json:"-" db:"password" validate:"required,min=6"` // hidden in JSON
	Role      UserRole  `json:"role" db:"role" validate:"required"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Achievement represents achievement entity
type Achievement struct {
	ID              string            `json:"id" db:"id"`
	MahasiswaID     string            `json:"mahasiswa_id" db:"mahasiswa_id" validate:"required"`
	Title           string            `json:"title" db:"title" validate:"required"`
	Description     string            `json:"description" db:"description" validate:"required"`
	Category        string            `json:"category" db:"category" validate:"required"`
	AchievementDate time.Time         `json:"achievement_date" db:"achievement_date" validate:"required"`
	Status          AchievementStatus `json:"status" db:"status"`
	VerifiedBy      *string           `json:"verified_by,omitempty" db:"verified_by"`
	VerifiedAt      *time.Time        `json:"verified_at,omitempty" db:"verified_at"`
	RejectionReason *string           `json:"rejection_reason,omitempty" db:"rejection_reason"`
	IsDeleted       bool              `json:"is_deleted" db:"is_deleted"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`
}

// AchievementFile represents file attachment for achievement (stored in MongoDB)
type AchievementFile struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	AchievementID string    `json:"achievement_id" bson:"achievement_id"`
	FileName      string    `json:"file_name" bson:"file_name"`
	FilePath      string    `json:"file_path" bson:"file_path"`
	FileSize      int64     `json:"file_size" bson:"file_size"`
	MimeType      string    `json:"mime_type" bson:"mime_type"`
	UploadedAt    time.Time `json:"uploaded_at" bson:"uploaded_at"`
}

// LoginRequest represents login request payload
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// CreateAchievementRequest represents create achievement request
type CreateAchievementRequest struct {
	Title           string    `json:"title" validate:"required"`
	Description     string    `json:"description" validate:"required"`
	Category        string    `json:"category" validate:"required"`
	AchievementDate time.Time `json:"achievement_date" validate:"required"`
}

// UpdateAchievementRequest represents update achievement request
type UpdateAchievementRequest struct {
	Title           *string    `json:"title,omitempty"`
	Description     *string    `json:"description,omitempty"`
	Category        *string    `json:"category,omitempty"`
	AchievementDate *time.Time `json:"achievement_date,omitempty"`
}

// RejectAchievementRequest represents reject achievement request
type RejectAchievementRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// APIResponse represents standard API response format
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

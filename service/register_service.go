package service

import (
	"database/sql"
	"errors"
	"prestasi-mahasiswa/utils"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type RegisterService struct {
	DB *sql.DB
}

func NewRegisterService(db *sql.DB) *RegisterService {
	return &RegisterService{DB: db}
}

func (s *RegisterService) CreateUser(nim, name, email, password, role string) error {
	// Validate required fields
	if err := s.ValidateUserData(nim, name, email, password, role); err != nil {
		return err
	}

	// Check if email already exists
	exists, err := s.checkEmailExists(email)
	if err != nil {
		return errors.New("database error: " + err.Error())
	}
	if exists {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return errors.New("failed to hash password: " + err.Error())
	}

	// Generate UUID for new user
	userID := uuid.New().String()

	// Insert user into database
	query := `
		INSERT INTO users (id, nim, name, email, password, role, is_active, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	now := time.Now()
	_, err = s.DB.Exec(query,
		userID,
		s.parseNIM(nim, role),
		name,
		email,
		hashedPassword,
		role,
		true, // is_active
		now,  // created_at
		now,  // updated_at
	)

	if err != nil {
		return errors.New("failed to create user: " + err.Error())
	}

	return nil
}

func (s *RegisterService) ValidateUserData(nim, name, email, password, role string) error {
	// Validate required fields
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(email) == "" {
		return errors.New("email is required")
	}
	if strings.TrimSpace(password) == "" {
		return errors.New("password is required")
	}

	// Validate password length
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	// Validate role
	validRoles := []string{"mahasiswa", "dosen_wali", "admin"}
	validRole := false
	for _, validRoleItem := range validRoles {
		if role == validRoleItem {
			validRole = true
			break
		}
	}
	if !validRole {
		return errors.New("invalid role. Must be mahasiswa, dosen_wali, or admin")
	}

	// Validate NIM for mahasiswa
	if role == "mahasiswa" && strings.TrimSpace(nim) == "" {
		return errors.New("NIM is required for mahasiswa")
	}

	return nil
}

func (s *RegisterService) checkEmailExists(email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = $1`
	var count int
	err := s.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SoftDeleteUser implements soft delete with status tracking
func (s *RegisterService) SoftDeleteUser(userID string, deletedBy string, reason string) error {
	query := `
		UPDATE users 
		SET is_deleted = true, is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND is_deleted = false
	`

	result, err := s.DB.Exec(query, userID)
	if err != nil {
		return errors.New("failed to delete user: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to check delete result: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("user not found or already deleted")
	}

	// TODO: Log deletion activity with reason and deleted_by
	// This would typically go to an audit_logs table

	return nil
}

func (s *RegisterService) parseNIM(nim, role string) *string {
	if role == "mahasiswa" && strings.TrimSpace(nim) != "" {
		return &nim
	}
	return nil
}

package service

import (
	"database/sql"
	"errors"
	"prestasi-mahasiswa/utils"

	_ "github.com/lib/pq"
)

type LoginService struct {
	DB      *sql.DB
	JWTUtil *utils.JWTUtil
}

func NewLoginService(db *sql.DB, jwtSecret string, jwtExpireHours int) *LoginService {
	return &LoginService{
		DB:      db,
		JWTUtil: utils.NewJWTUtil(jwtSecret, jwtExpireHours),
	}
}

func (s *LoginService) AuthenticateUser(email, password string) (*UserData, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	// Query user from database (handle soft delete with fallback)
	query := `SELECT id, nim, name, email, password, role, is_active FROM users WHERE email = $1 AND is_active = true`
	var user UserData
	var hashedPassword string
	var isActive bool

	err := s.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.NIM,
		&user.Name,
		&user.Email,
		&hashedPassword,
		&user.Role,
		&isActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid email or password or account inactive")
		}
		return nil, errors.New("database error: " + err.Error())
	}

	// Verify password
	err = utils.ComparePassword(hashedPassword, password)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Set user status and flags
	user.IsActive = isActive
	user.Status = "active" // Only active users can login with current query

	return &user, nil
}

func (s *LoginService) GenerateToken(user *UserData) (string, error) {
	return s.JWTUtil.GenerateToken(user.ID, user.Email, user.Role)
}

// GetUserInfo retrieves user information by ID
func (s *LoginService) GetUserInfo(userID string) (*UserData, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	query := `SELECT id, nim, name, email, role, is_active FROM users WHERE id = $1 AND is_active = true`
	var user UserData

	err := s.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.NIM,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("database error: " + err.Error())
	}

	user.Status = "active"
	return &user, nil
}

type UserData struct {
	ID       string  `json:"id"`
	NIM      *string `json:"nim,omitempty"`
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	Role     string  `json:"role"`
	IsActive bool    `json:"is_active"`
	Status   string  `json:"status"` // active, inactive, deleted
}

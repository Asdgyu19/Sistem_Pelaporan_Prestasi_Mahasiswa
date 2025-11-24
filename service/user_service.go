package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB *sql.DB
}

type User struct {
	ID          string                 `json:"id"`
	NIM         *string                `json:"nim,omitempty"`
	NIP         *string                `json:"nip,omitempty"`
	Name        string                 `json:"name"`
	Email       string                 `json:"email"`
	Role        string                 `json:"role"`
	AdvisorID   *string                `json:"advisor_id,omitempty"`
	AdvisorName *string                `json:"advisor_name,omitempty"`
	ProfileData map[string]interface{} `json:"profile_data,omitempty"`
	IsActive    bool                   `json:"is_active"`
	IsDeleted   bool                   `json:"is_deleted"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type CreateUserRequest struct {
	NIM         *string                `json:"nim,omitempty"`
	NIP         *string                `json:"nip,omitempty"`
	Name        string                 `json:"name"`
	Email       string                 `json:"email"`
	Password    string                 `json:"password"`
	Role        string                 `json:"role"`
	AdvisorID   *string                `json:"advisor_id,omitempty"`
	ProfileData map[string]interface{} `json:"profile_data,omitempty"`
}

type UpdateUserRequest struct {
	NIM         *string                `json:"nim,omitempty"`
	NIP         *string                `json:"nip,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Email       *string                `json:"email,omitempty"`
	Password    *string                `json:"password,omitempty"`
	Role        *string                `json:"role,omitempty"`
	AdvisorID   *string                `json:"advisor_id,omitempty"`
	ProfileData map[string]interface{} `json:"profile_data,omitempty"`
	IsActive    *bool                  `json:"is_active,omitempty"`
}

type AssignAdvisorRequest struct {
	MahasiswaID string `json:"mahasiswa_id"`
	AdvisorID   string `json:"advisor_id"`
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

// Get all users with optional filters
func (s *UserService) GetAllUsers(role, search string, includeDeleted bool) ([]User, error) {
	// Check which columns exist in the database
	checkColumnsQuery := `
		SELECT 
			COUNT(CASE WHEN column_name = 'advisor_id' THEN 1 END) > 0 as has_advisor,
			COUNT(CASE WHEN column_name = 'profile_data' THEN 1 END) > 0 as has_profile,
			COUNT(CASE WHEN column_name = 'is_deleted' THEN 1 END) > 0 as has_deleted
		FROM information_schema.columns 
		WHERE table_name = 'users' AND column_name IN ('advisor_id', 'profile_data', 'is_deleted')
	`

	var hasAdvisor, hasProfile, hasDeleted bool
	err := s.DB.QueryRow(checkColumnsQuery).Scan(&hasAdvisor, &hasProfile, &hasDeleted)
	if err != nil {
		return nil, fmt.Errorf("failed to check database schema: %v", err)
	}

	// Build query based on available columns
	selectClause := "u.id, u.nim, u.nip, u.name, u.email, u.role"

	if hasAdvisor {
		selectClause += ", u.advisor_id, a.name as advisor_name"
	} else {
		selectClause += ", NULL as advisor_id, NULL as advisor_name"
	}

	if hasProfile {
		selectClause += ", u.profile_data"
	} else {
		selectClause += ", NULL as profile_data"
	}

	selectClause += ", u.is_active"

	if hasDeleted {
		selectClause += ", u.is_deleted"
	} else {
		selectClause += ", FALSE as is_deleted"
	}

	selectClause += ", u.created_at, u.updated_at"

	var fromClause string
	if hasAdvisor {
		fromClause = "FROM users u LEFT JOIN users a ON u.advisor_id = a.id"
	} else {
		fromClause = "FROM users u"
	}

	query := fmt.Sprintf("SELECT %s %s WHERE 1=1", selectClause, fromClause)
	args := []interface{}{}
	argCount := 0

	if !includeDeleted && hasDeleted {
		query += " AND u.is_deleted = FALSE"
	}

	if role != "" {
		argCount++
		query += fmt.Sprintf(" AND u.role = $%d", argCount)
		args = append(args, role)
	}

	if search != "" {
		argCount++
		query += fmt.Sprintf(" AND (u.name ILIKE $%d OR u.email ILIKE $%d OR u.nim ILIKE $%d OR u.nip ILIKE $%d)",
			argCount, argCount, argCount, argCount)
		args = append(args, "%"+search+"%")
	}

	query += " ORDER BY u.created_at DESC"

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var profileDataJSON []byte

		err := rows.Scan(
			&user.ID, &user.NIM, &user.NIP, &user.Name, &user.Email, &user.Role,
			&user.AdvisorID, &user.AdvisorName, &profileDataJSON,
			&user.IsActive, &user.IsDeleted, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}

		// Parse profile data JSON if it exists
		if len(profileDataJSON) > 0 {
			json.Unmarshal(profileDataJSON, &user.ProfileData)
		}

		users = append(users, user)
	}

	return users, nil
}

// Get user by ID
func (s *UserService) GetUserByID(userID string) (*User, error) {
	// Check which columns exist
	checkColumnsQuery := `
		SELECT 
			COUNT(CASE WHEN column_name = 'advisor_id' THEN 1 END) > 0 as has_advisor,
			COUNT(CASE WHEN column_name = 'profile_data' THEN 1 END) > 0 as has_profile,
			COUNT(CASE WHEN column_name = 'is_deleted' THEN 1 END) > 0 as has_deleted
		FROM information_schema.columns 
		WHERE table_name = 'users' AND column_name IN ('advisor_id', 'profile_data', 'is_deleted')
	`

	var hasAdvisor, hasProfile, hasDeleted bool
	err := s.DB.QueryRow(checkColumnsQuery).Scan(&hasAdvisor, &hasProfile, &hasDeleted)
	if err != nil {
		return nil, fmt.Errorf("failed to check database schema: %v", err)
	}

	// Build query based on available columns
	selectClause := "u.id, u.nim, u.nip, u.name, u.email, u.role"

	if hasAdvisor {
		selectClause += ", u.advisor_id, a.name as advisor_name"
	} else {
		selectClause += ", NULL as advisor_id, NULL as advisor_name"
	}

	if hasProfile {
		selectClause += ", u.profile_data"
	} else {
		selectClause += ", NULL as profile_data"
	}

	selectClause += ", u.is_active"

	if hasDeleted {
		selectClause += ", u.is_deleted"
	} else {
		selectClause += ", FALSE as is_deleted"
	}

	selectClause += ", u.created_at, u.updated_at"

	var fromClause, whereClause string
	if hasAdvisor {
		fromClause = "FROM users u LEFT JOIN users a ON u.advisor_id = a.id"
	} else {
		fromClause = "FROM users u"
	}

	if hasDeleted {
		whereClause = "WHERE u.id = $1 AND u.is_deleted = FALSE"
	} else {
		whereClause = "WHERE u.id = $1"
	}

	query := fmt.Sprintf("SELECT %s %s %s", selectClause, fromClause, whereClause)

	var user User
	var profileDataJSON []byte

	err = s.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.NIM, &user.NIP, &user.Name, &user.Email, &user.Role,
		&user.AdvisorID, &user.AdvisorName, &profileDataJSON,
		&user.IsActive, &user.IsDeleted, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Parse profile data JSON
	if len(profileDataJSON) > 0 {
		json.Unmarshal(profileDataJSON, &user.ProfileData)
	}

	return &user, nil
}

// Create new user
func (s *UserService) CreateUser(req CreateUserRequest) (*User, error) {
	// Validate required fields
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" {
		return nil, errors.New("missing required fields: name, email, password, role")
	}

	// Validate role
	if req.Role != "mahasiswa" && req.Role != "dosen_wali" && req.Role != "admin" {
		return nil, errors.New("invalid role. Must be: mahasiswa, dosen_wali, or admin")
	}

	// Validate role-specific fields
	if req.Role == "mahasiswa" && req.NIM == nil {
		return nil, errors.New("NIM is required for mahasiswa")
	}
	if (req.Role == "dosen_wali" || req.Role == "admin") && req.NIP == nil {
		return nil, errors.New("NIP is required for dosen_wali and admin")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Convert profile data to JSON
	var profileDataJSON []byte
	if req.ProfileData != nil {
		profileDataJSON, _ = json.Marshal(req.ProfileData)
	}

	userID := uuid.New().String()

	// Check which columns exist
	checkColumnsQuery := `
		SELECT 
			COUNT(CASE WHEN column_name = 'advisor_id' THEN 1 END) > 0 as has_advisor,
			COUNT(CASE WHEN column_name = 'profile_data' THEN 1 END) > 0 as has_profile,
			COUNT(CASE WHEN column_name = 'is_deleted' THEN 1 END) > 0 as has_deleted
		FROM information_schema.columns 
		WHERE table_name = 'users' AND column_name IN ('advisor_id', 'profile_data', 'is_deleted')
	`

	var hasAdvisor, hasProfile, hasDeleted bool
	err = s.DB.QueryRow(checkColumnsQuery).Scan(&hasAdvisor, &hasProfile, &hasDeleted)
	if err != nil {
		return nil, fmt.Errorf("failed to check database schema: %v", err)
	}

	// Build INSERT query based on available columns
	columns := "id, nim, nip, name, email, password, role, is_active"
	values := "$1, $2, $3, $4, $5, $6, $7, TRUE"
	args := []interface{}{userID, req.NIM, req.NIP, req.Name, req.Email, string(hashedPassword), req.Role}
	argCount := 7

	if hasAdvisor {
		argCount++
		columns += ", advisor_id"
		values += fmt.Sprintf(", $%d", argCount)
		args = append(args, req.AdvisorID)
	}
	

	if hasProfile {
		argCount++
		columns += ", profile_data"
		values += fmt.Sprintf(", $%d", argCount)
		args = append(args, profileDataJSON)
	}

	if hasDeleted {
		columns += ", is_deleted"
		values += ", FALSE"
	}

	query := fmt.Sprintf(`
		INSERT INTO users (%s)
		VALUES (%s)
		RETURNING created_at, updated_at
	`, columns, values)

	var createdAt, updatedAt time.Time
	err = s.DB.QueryRow(query, args...).Scan(&createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Return created user
	return s.GetUserByID(userID)
}

// Update user
func (s *UserService) UpdateUser(userID string, req UpdateUserRequest) (*User, error) {
	// Check if user exists
	existingUser, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argCount := 0

	if req.NIM != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("nim = $%d", argCount))
		args = append(args, req.NIM)
	}

	if req.NIP != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("nip = $%d", argCount))
		args = append(args, req.NIP)
	}

	if req.Name != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("name = $%d", argCount))
		args = append(args, *req.Name)
	}

	if req.Email != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("email = $%d", argCount))
		args = append(args, *req.Email)
	}

	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
		argCount++
		setParts = append(setParts, fmt.Sprintf("password = $%d", argCount))
		args = append(args, string(hashedPassword))
	}

	if req.Role != nil {
		if *req.Role != "mahasiswa" && *req.Role != "dosen_wali" && *req.Role != "admin" {
			return nil, errors.New("invalid role. Must be: mahasiswa, dosen_wali, or admin")
		}
		argCount++
		setParts = append(setParts, fmt.Sprintf("role = $%d", argCount))
		args = append(args, *req.Role)
	}

	if req.AdvisorID != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("advisor_id = $%d", argCount))
		args = append(args, req.AdvisorID)
	}

	if req.ProfileData != nil {
		profileDataJSON, _ := json.Marshal(req.ProfileData)
		argCount++
		setParts = append(setParts, fmt.Sprintf("profile_data = $%d", argCount))
		args = append(args, profileDataJSON)
	}

	if req.IsActive != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argCount))
		args = append(args, *req.IsActive)
	}

	if len(setParts) == 0 {
		return existingUser, nil // No changes
	}

	// Add updated_at
	argCount++
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())

	// Add WHERE clause
	argCount++
	args = append(args, userID)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d",
		strings.Join(setParts, ", "), argCount)

	_, err = s.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return s.GetUserByID(userID)
}

// Soft delete user
func (s *UserService) DeleteUser(userID string) error {
	// Check if is_deleted column exists
	var hasDeleted bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_name = 'users' AND column_name = 'is_deleted'
		)
	`
	err := s.DB.QueryRow(checkQuery).Scan(&hasDeleted)
	if err != nil {
		return fmt.Errorf("failed to check database schema: %v", err)
	}

	var query string
	if hasDeleted {
		query = "UPDATE users SET is_deleted = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND is_deleted = FALSE"
	} else {
		// If no is_deleted column, we can't soft delete, so return error or use hard delete
		return errors.New("soft delete not supported - is_deleted column missing. Please add the column to database")
	}

	result, err := s.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found or already deleted")
	}

	return nil
} // Assign advisor to mahasiswa
func (s *UserService) AssignAdvisor(req AssignAdvisorRequest) error {
	// Verify mahasiswa exists and is mahasiswa role
	mahasiswa, err := s.GetUserByID(req.MahasiswaID)
	if err != nil {
		return fmt.Errorf("mahasiswa not found: %v", err)
	}
	if mahasiswa.Role != "mahasiswa" {
		return errors.New("user is not a mahasiswa")
	}

	// Verify advisor exists and is dosen_wali
	advisor, err := s.GetUserByID(req.AdvisorID)
	if err != nil {
		return fmt.Errorf("advisor not found: %v", err)
	}
	if advisor.Role != "dosen_wali" {
		return errors.New("advisor must be a dosen_wali")
	}

	// Update mahasiswa with advisor
	query := "UPDATE users SET advisor_id = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2"
	_, err = s.DB.Exec(query, req.AdvisorID, req.MahasiswaID)
	if err != nil {
		return fmt.Errorf("failed to assign advisor: %v", err)
	}

	return nil
}

// Get students by advisor
func (s *UserService) GetStudentsByAdvisor(advisorID string) ([]User, error) {
	query := `
		SELECT 
			u.id, u.nim, u.nip, u.name, u.email, u.role, 
			u.advisor_id, a.name as advisor_name, u.profile_data,
			u.is_active, u.is_deleted, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN users a ON u.advisor_id = a.id
		WHERE u.advisor_id = $1 AND u.role = 'mahasiswa' AND u.is_deleted = FALSE
		ORDER BY u.name ASC
	`

	rows, err := s.DB.Query(query, advisorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get students: %v", err)
	}
	defer rows.Close()

	var students []User
	for rows.Next() {
		var student User
		var profileDataJSON []byte

		err := rows.Scan(
			&student.ID, &student.NIM, &student.NIP, &student.Name, &student.Email, &student.Role,
			&student.AdvisorID, &student.AdvisorName, &profileDataJSON,
			&student.IsActive, &student.IsDeleted, &student.CreatedAt, &student.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student: %v", err)
		}

		// Parse profile data JSON
		if len(profileDataJSON) > 0 {
			json.Unmarshal(profileDataJSON, &student.ProfileData)
		}

		students = append(students, student)
	}

	return students, nil
}

// Get available advisors (dosen_wali)
func (s *UserService) GetAvailableAdvisors() ([]User, error) {
	return s.GetAllUsers("dosen_wali", "", false)
}

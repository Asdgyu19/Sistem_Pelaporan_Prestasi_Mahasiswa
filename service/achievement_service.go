package service

import (
	"database/sql"
	"errors"
	"prestasi-mahasiswa/database"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type AchievementService struct {
	DB      *sql.DB
	MongoDB *database.MongoDB
}

func NewAchievementService(db *sql.DB, mongodb *database.MongoDB) *AchievementService {
	return &AchievementService{
		DB:      db,
		MongoDB: mongodb,
	}
}

// GetAllAchievements retrieves all achievements (for dosen/admin)
func (s *AchievementService) GetAllAchievements() ([]Achievement, error) {
	query := `
		SELECT a.id, a.mahasiswa_id, u.name as mahasiswa_name, a.title, a.description, 
		       a.category, a.achievement_date, a.status, a.verified_by, a.verified_at,
		       a.rejection_reason, a.created_at, a.updated_at
		FROM achievements a
		JOIN users u ON a.mahasiswa_id = u.id
		WHERE a.is_deleted = false OR a.is_deleted IS NULL
		ORDER BY a.created_at DESC
	`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, errors.New("failed to fetch achievements: " + err.Error())
	}
	defer rows.Close()

	var achievements []Achievement
	for rows.Next() {
		var a Achievement
		var verifiedAt, rejectionReason sql.NullString
		var verifiedBy sql.NullString

		err := rows.Scan(
			&a.ID, &a.MahasiswaID, &a.MahasiswaName, &a.Title, &a.Description,
			&a.Category, &a.AchievementDate, &a.Status, &verifiedBy, &verifiedAt,
			&rejectionReason, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, errors.New("failed to scan achievement: " + err.Error())
		}

		if verifiedBy.Valid {
			a.VerifiedBy = &verifiedBy.String
		}
		if verifiedAt.Valid {
			a.VerifiedAt = &verifiedAt.String
		}
		if rejectionReason.Valid {
			a.RejectionReason = &rejectionReason.String
		}

		achievements = append(achievements, a)
	}

	return achievements, nil
}

// GetAchievementsByMahasiswa retrieves achievements for specific student
func (s *AchievementService) GetAchievementsByMahasiswa(mahasiswaID string) ([]Achievement, error) {
	query := `
		SELECT a.id, a.mahasiswa_id, u.name as mahasiswa_name, a.title, a.description,
		       a.category, a.achievement_date, a.status, a.verified_by, a.verified_at,
		       a.rejection_reason, a.created_at, a.updated_at
		FROM achievements a
		JOIN users u ON a.mahasiswa_id = u.id
		WHERE a.mahasiswa_id = $1 AND (a.is_deleted = false OR a.is_deleted IS NULL)
		ORDER BY a.created_at DESC
	`

	rows, err := s.DB.Query(query, mahasiswaID)
	if err != nil {
		return nil, errors.New("failed to fetch achievements: " + err.Error())
	}
	defer rows.Close()

	var achievements []Achievement
	for rows.Next() {
		var a Achievement
		var verifiedAt, rejectionReason sql.NullString
		var verifiedBy sql.NullString

		err := rows.Scan(
			&a.ID, &a.MahasiswaID, &a.MahasiswaName, &a.Title, &a.Description,
			&a.Category, &a.AchievementDate, &a.Status, &verifiedBy, &verifiedAt,
			&rejectionReason, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, errors.New("failed to scan achievement: " + err.Error())
		}

		if verifiedBy.Valid {
			a.VerifiedBy = &verifiedBy.String
		}
		if verifiedAt.Valid {
			a.VerifiedAt = &verifiedAt.String
		}
		if rejectionReason.Valid {
			a.RejectionReason = &rejectionReason.String
		}

		achievements = append(achievements, a)
	}

	return achievements, nil
}

// CreateAchievement creates new achievement
func (s *AchievementService) CreateAchievement(achievement Achievement) (*Achievement, error) {
	// Validate required fields
	if achievement.Title == "" || achievement.Description == "" || achievement.Category == "" {
		return nil, errors.New("title, description, and category are required")
	}

	if achievement.MahasiswaID == "" {
		return nil, errors.New("mahasiswa_id is required")
	}

	// Generate UUID for new achievement
	achievementID := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO achievements (id, mahasiswa_id, title, description, category, 
		                         achievement_date, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	var createdID string
	err := s.DB.QueryRow(query,
		achievementID,
		achievement.MahasiswaID,
		achievement.Title,
		achievement.Description,
		achievement.Category,
		achievement.AchievementDate,
		"draft", // Default status
		now,
		now,
	).Scan(&createdID)

	if err != nil {
		return nil, errors.New("failed to create achievement: " + err.Error())
	}

	// Return created achievement
	achievement.ID = createdID
	achievement.Status = "draft"
	achievement.CreatedAt = now.Format(time.RFC3339)
	achievement.UpdatedAt = now.Format(time.RFC3339)

	return &achievement, nil
}

// GetAchievementByID retrieves single achievement by ID
func (s *AchievementService) GetAchievementByID(id string) (*Achievement, error) {
	query := `
		SELECT a.id, a.mahasiswa_id, u.name as mahasiswa_name, a.title, a.description,
		       a.category, a.achievement_date, a.status, a.verified_by, a.verified_at,
		       a.rejection_reason, a.created_at, a.updated_at
		FROM achievements a
		JOIN users u ON a.mahasiswa_id = u.id
		WHERE a.id = $1 AND (a.is_deleted = false OR a.is_deleted IS NULL)
	`

	var a Achievement
	var verifiedAt, rejectionReason sql.NullString
	var verifiedBy sql.NullString

	err := s.DB.QueryRow(query, id).Scan(
		&a.ID, &a.MahasiswaID, &a.MahasiswaName, &a.Title, &a.Description,
		&a.Category, &a.AchievementDate, &a.Status, &verifiedBy, &verifiedAt,
		&rejectionReason, &a.CreatedAt, &a.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("achievement not found")
		}
		return nil, errors.New("failed to fetch achievement: " + err.Error())
	}

	if verifiedBy.Valid {
		a.VerifiedBy = &verifiedBy.String
	}
	if verifiedAt.Valid {
		a.VerifiedAt = &verifiedAt.String
	}
	if rejectionReason.Valid {
		a.RejectionReason = &rejectionReason.String
	}

	return &a, nil
}

// UpdateAchievement updates existing achievement
func (s *AchievementService) UpdateAchievement(id string, achievement Achievement) error {
	query := `
		UPDATE achievements 
		SET title = $1, description = $2, category = $3, achievement_date = $4, updated_at = $5
		WHERE id = $6 AND status = 'draft' AND (is_deleted = false OR is_deleted IS NULL)
	`

	result, err := s.DB.Exec(query,
		achievement.Title,
		achievement.Description,
		achievement.Category,
		achievement.AchievementDate,
		time.Now(),
		id,
	)

	if err != nil {
		return errors.New("failed to update achievement: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to check update result: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("achievement not found or cannot be updated (only draft status can be updated)")
	}

	return nil
}

// DeleteAchievement soft deletes achievement
func (s *AchievementService) DeleteAchievement(id string) error {
	query := `
		UPDATE achievements 
		SET is_deleted = true, updated_at = $1
		WHERE id = $2 AND status = 'draft'
	`

	result, err := s.DB.Exec(query, time.Now(), id)
	if err != nil {
		return errors.New("failed to delete achievement: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to check delete result: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("achievement not found or cannot be deleted (only draft status can be deleted)")
	}

	return nil
}

// SubmitAchievement submits achievement for verification
func (s *AchievementService) SubmitAchievement(id string) error {
	query := `
		UPDATE achievements 
		SET status = 'submitted', updated_at = $1
		WHERE id = $2 AND status = 'draft' AND (is_deleted = false OR is_deleted IS NULL)
	`

	result, err := s.DB.Exec(query, time.Now(), id)
	if err != nil {
		return errors.New("failed to submit achievement: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to check submit result: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("achievement not found or cannot be submitted")
	}

	return nil
}

// VerifyAchievement verifies achievement (for dosen/admin)
func (s *AchievementService) VerifyAchievement(id string, verifierID string) error {
	query := `
		UPDATE achievements 
		SET status = 'verified', verified_by = $1, verified_at = $2, updated_at = $3
		WHERE id = $4 AND status = 'submitted' AND (is_deleted = false OR is_deleted IS NULL)
	`

	now := time.Now()
	result, err := s.DB.Exec(query, verifierID, now, now, id)
	if err != nil {
		return errors.New("failed to verify achievement: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to check verify result: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("achievement not found or cannot be verified")
	}

	return nil
}

// RejectAchievement rejects achievement with reason (for dosen/admin)
func (s *AchievementService) RejectAchievement(id string, verifierID string, reason string) error {
	query := `
		UPDATE achievements 
		SET status = 'rejected', verified_by = $1, verified_at = $2, rejection_reason = $3, updated_at = $4
		WHERE id = $5 AND status = 'submitted' AND (is_deleted = false OR is_deleted IS NULL)
	`

	now := time.Now()
	result, err := s.DB.Exec(query, verifierID, now, reason, now, id)
	if err != nil {
		return errors.New("failed to reject achievement: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to check reject result: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("achievement not found or cannot be rejected")
	}

	return nil
}

type Achievement struct {
	ID              string    `json:"id"`
	MahasiswaID     string    `json:"mahasiswa_id"`
	MahasiswaName   string    `json:"mahasiswa_name,omitempty"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	AchievementDate time.Time `json:"achievement_date"`
	Status          string    `json:"status"` // draft, submitted, verified, rejected
	VerifiedBy      *string   `json:"verified_by,omitempty"`
	VerifiedAt      *string   `json:"verified_at,omitempty"`
	RejectionReason *string   `json:"rejection_reason,omitempty"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
}

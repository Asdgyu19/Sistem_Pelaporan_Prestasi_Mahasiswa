package service

import (
	"database/sql"
	"fmt"
	"time"
)

type ReportService struct {
	DB *sql.DB
}

func NewReportService(db *sql.DB) *ReportService {
	return &ReportService{DB: db}
}

type AchievementStatistics struct {
	TotalAchievements    int64  `json:"total_achievements"`
	VerifiedAchievements int64  `json:"verified_achievements"`
	PendingAchievements  int64  `json:"pending_achievements"`
	RejectedAchievements int64  `json:"rejected_achievements"`
	VerificationRate     string `json:"verification_rate"`
	AvgVerificationTime  string `json:"avg_verification_time"`
}

type SystemStatistics struct {
	TotalUsers     int64                 `json:"total_users"`
	TotalMahasiswa int64                 `json:"total_mahasiswa"`
	TotalDosen     int64                 `json:"total_dosen"`
	TotalAdmins    int64                 `json:"total_admins"`
	Achievements   AchievementStatistics `json:"achievements"`
	TopCategories  []CategoryCount       `json:"top_categories"`
}

type CategoryCount struct {
	Category string `json:"category"`
	Count    int64  `json:"count"`
}

type StudentReport struct {
	StudentID            string           `json:"student_id"`
	StudentName          string           `json:"student_name"`
	Email                string           `json:"email"`
	NIM                  *string          `json:"nim"`
	AdvisorID            *string          `json:"advisor_id"`
	AdvisorName          *string          `json:"advisor_name"`
	TotalAchievements    int64            `json:"total_achievements"`
	VerifiedAchievements int64            `json:"verified_achievements"`
	PendingAchievements  int64            `json:"pending_achievements"`
	RejectedAchievements int64            `json:"rejected_achievements"`
	VerificationRate     string           `json:"verification_rate"`
	AchievementsByStatus map[string]int64 `json:"achievements_by_status"`
	AchievementsByType   map[string]int64 `json:"achievements_by_category"`
}

type AchievementHistory struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	ChangedBy    string    `json:"changed_by"`
	ChangedAt    time.Time `json:"changed_at"`
	Reason       *string   `json:"reason"`
	VerifiedName *string   `json:"verified_name,omitempty"`
}

// GetSystemStatistics retrieves overall system statistics
func (rs *ReportService) GetSystemStatistics() (*SystemStatistics, error) {
	stats := &SystemStatistics{}

	// Get user counts
	userQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN role = 'mahasiswa' THEN 1 END) as mahasiswa,
			COUNT(CASE WHEN role = 'dosen_wali' THEN 1 END) as dosen,
			COUNT(CASE WHEN role = 'admin' THEN 1 END) as admins
		FROM users WHERE is_active = true
	`

	err := rs.DB.QueryRow(userQuery).Scan(
		&stats.TotalUsers,
		&stats.TotalMahasiswa,
		&stats.TotalDosen,
		&stats.TotalAdmins,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get user statistics: %w", err)
	}

	// Get achievement statistics
	achQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'verified' THEN 1 END) as verified,
			COUNT(CASE WHEN status = 'submitted' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'rejected' THEN 1 END) as rejected
		FROM achievements WHERE is_deleted = false
	`

	var total, verified, pending, rejected int64
	err = rs.DB.QueryRow(achQuery).Scan(&total, &verified, &pending, &rejected)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get achievement statistics: %w", err)
	}

	stats.Achievements.TotalAchievements = total
	stats.Achievements.VerifiedAchievements = verified
	stats.Achievements.PendingAchievements = pending
	stats.Achievements.RejectedAchievements = rejected

	// Calculate verification rate
	if total > 0 {
		rate := float64(verified) / float64(total) * 100
		stats.Achievements.VerificationRate = fmt.Sprintf("%.2f%%", rate)
	} else {
		stats.Achievements.VerificationRate = "0%"
	}

	// Get top categories
	categoryQuery := `
		SELECT category, COUNT(*) as count
		FROM achievements
		WHERE is_deleted = false
		GROUP BY category
		ORDER BY count DESC
		LIMIT 5
	`

	rows, err := rs.DB.Query(categoryQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get top categories: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cat CategoryCount
		if err := rows.Scan(&cat.Category, &cat.Count); err != nil {
			continue
		}
		stats.TopCategories = append(stats.TopCategories, cat)
	}

	return stats, nil
}

// GetStudentReport retrieves comprehensive report for a specific student
func (rs *ReportService) GetStudentReport(studentID string) (*StudentReport, error) {
	report := &StudentReport{StudentID: studentID}

	// Get student info
	studentQuery := `
		SELECT id, name, email, nim, advisor_id FROM users 
		WHERE id = $1 AND role = 'mahasiswa' AND is_active = true
	`

	err := rs.DB.QueryRow(studentQuery, studentID).Scan(
		&report.StudentID,
		&report.StudentName,
		&report.Email,
		&report.NIM,
		&report.AdvisorID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student not found")
		}
		return nil, fmt.Errorf("failed to get student info: %w", err)
	}

	// Get advisor name if exists
	if report.AdvisorID != nil {
		advisorQuery := `SELECT name FROM users WHERE id = $1`
		rs.DB.QueryRow(advisorQuery, *report.AdvisorID).Scan(&report.AdvisorName)
	}

	// Get achievement statistics
	achQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'verified' THEN 1 END) as verified,
			COUNT(CASE WHEN status = 'submitted' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'rejected' THEN 1 END) as rejected
		FROM achievements
		WHERE mahasiswa_id = $1 AND is_deleted = false
	`

	err = rs.DB.QueryRow(achQuery, studentID).Scan(
		&report.TotalAchievements,
		&report.VerifiedAchievements,
		&report.PendingAchievements,
		&report.RejectedAchievements,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get achievement statistics: %w", err)
	}

	// Calculate verification rate
	if report.TotalAchievements > 0 {
		rate := float64(report.VerifiedAchievements) / float64(report.TotalAchievements) * 100
		report.VerificationRate = fmt.Sprintf("%.2f%%", rate)
	} else {
		report.VerificationRate = "0%"
	}

	// Get achievements by status
	statusQuery := `
		SELECT status, COUNT(*) as count
		FROM achievements
		WHERE mahasiswa_id = $1 AND is_deleted = false
		GROUP BY status
	`

	report.AchievementsByStatus = make(map[string]int64)
	rows, err := rs.DB.Query(statusQuery, studentID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var status string
			var count int64
			if err := rows.Scan(&status, &count); err != nil {
				continue
			}
			report.AchievementsByStatus[status] = count
		}
	}

	// Get achievements by category
	categoryQuery := `
		SELECT category, COUNT(*) as count
		FROM achievements
		WHERE mahasiswa_id = $1 AND is_deleted = false
		GROUP BY category
	`

	report.AchievementsByType = make(map[string]int64)
	rows, err = rs.DB.Query(categoryQuery, studentID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var category string
			var count int64
			if err := rows.Scan(&category, &count); err != nil {
				continue
			}
			report.AchievementsByType[category] = count
		}
	}

	return report, nil
}

// GetAchievementHistory retrieves status change history for an achievement
func (rs *ReportService) GetAchievementHistory(achievementID string) ([]AchievementHistory, error) {
	query := `
		SELECT 
			id,
			status,
			verified_by,
			verified_at,
			rejection_reason,
			u.name
		FROM achievements a
		LEFT JOIN users u ON a.verified_by = u.id
		WHERE a.id = $1
		ORDER BY a.verified_at DESC
	`

	rows, err := rs.DB.Query(query, achievementID)
	if err != nil {
		return nil, fmt.Errorf("failed to get achievement history: %w", err)
	}
	defer rows.Close()

	var history []AchievementHistory

	for rows.Next() {
		var h AchievementHistory
		var verifiedBy sql.NullString
		var verifiedAt sql.NullTime
		var rejectionReason sql.NullString
		var verifiedName sql.NullString

		err := rows.Scan(&h.ID, &h.Status, &verifiedBy, &verifiedAt, &rejectionReason, &verifiedName)
		if err != nil {
			continue
		}

		if verifiedBy.Valid {
			h.ChangedBy = verifiedBy.String
		}
		if verifiedAt.Valid {
			h.ChangedAt = verifiedAt.Time
		}
		if rejectionReason.Valid {
			h.Reason = &rejectionReason.String
		}
		if verifiedName.Valid {
			h.VerifiedName = &verifiedName.String
		}

		history = append(history, h)
	}

	return history, nil
}

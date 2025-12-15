package helper

import (
	"fmt"
	"io"
	"prestasi-mahasiswa/service"
	"time"

	"github.com/gin-gonic/gin"
)

type AchievementHelper struct {
	AchievementService *service.AchievementService
	FileService        *service.FileService
}

func NewAchievementHelper(achievementSvc *service.AchievementService, fileSvc *service.FileService) *AchievementHelper {
	return &AchievementHelper{
		AchievementService: achievementSvc,
		FileService:        fileSvc,
	}
}

// GetAchievements godoc
func (h *AchievementHelper) GetAchievements(c *gin.Context) {
	// Check user role from middleware
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(401, gin.H{"error": "User role not found"})
		return
	}

	userID, userIDExists := c.Get("user_id")
	if !userIDExists {
		c.JSON(401, gin.H{"error": "User ID not found"})
		return
	}

	var achievements []service.Achievement
	var err error

	// Mahasiswa only see their own achievements
	if userRole.(string) == "mahasiswa" {
		achievements, err = h.AchievementService.GetAchievementsByMahasiswa(userID.(string))
	} else {
		// Dosen/Admin see all achievements
		achievements, err = h.AchievementService.GetAllAchievements()
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Get achievements",
		"data":    achievements,
		"count":   len(achievements),
		"status":  "success",
	})
}


// @Router /achievements [post]
func (h *AchievementHelper) CreateAchievement(c *gin.Context) {
	// Get mahasiswa ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "User ID not found"})
		return
	}

	userRole, roleExists := c.Get("user_role")
	if !roleExists || userRole.(string) != "mahasiswa" {
		c.JSON(403, gin.H{"error": "Only mahasiswa can create achievements"})
		return
	}

	var req struct {
		Title           string    `json:"title" binding:"required"`
		Description     string    `json:"description" binding:"required"`
		Category        string    `json:"category" binding:"required"`
		AchievementDate time.Time `json:"achievement_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Create achievement object
	achievement := service.Achievement{
		MahasiswaID:     userID.(string),
		Title:           req.Title,
		Description:     req.Description,
		Category:        req.Category,
		AchievementDate: req.AchievementDate,
	}

	createdAchievement, err := h.AchievementService.CreateAchievement(achievement)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"message": "Achievement created successfully",
		"data":    createdAchievement,
		"status":  "success",
	})
}


func (h *AchievementHelper) GetAchievement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Achievement ID is required"})
		return
	}

	achievement, err := h.AchievementService.GetAchievementByID(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Check access rights
	userRole, _ := c.Get("user_role")
	userID, _ := c.Get("user_id")

	// Mahasiswa can only see their own achievements
	if userRole.(string) == "mahasiswa" && achievement.MahasiswaID != userID.(string) {
		c.JSON(403, gin.H{"error": "Access denied: You can only view your own achievements"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Get achievement by ID: " + id,
		"data":    achievement,
		"status":  "success",
	})
}

// UpdateAchievement godoc
func (h *AchievementHelper) UpdateAchievement(c *gin.Context) {
	id := c.Param("id")
	var achievement service.Achievement
	if err := c.ShouldBindJSON(&achievement); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.AchievementService.UpdateAchievement(id, achievement); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Achievement updated successfully",
		"status":  "success",
	})
}

// DeleteAchievement godoc
func (h *AchievementHelper) DeleteAchievement(c *gin.Context) {
	id := c.Param("id")
	if err := h.AchievementService.DeleteAchievement(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Achievement deleted successfully",
		"status":  "success",
	})
}

// UploadFile godoc
func (h *AchievementHelper) UploadFile(c *gin.Context) {
	achievementID := c.Param("id")
	if achievementID == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Achievement ID is required",
		})
		return
	}

	// Check if achievement exists and user has access
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"message": "User ID not found",
		})
		return
	}

	userRole, _ := c.Get("user_role")

	// Verify achievement exists and user can upload to it
	achievement, err := h.AchievementService.GetAchievementByID(achievementID)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Achievement not found",
			"error":   err.Error(),
		})
		return
	}

	// Check access rights - mahasiswa can only upload to their own achievements
	if userRole.(string) == "mahasiswa" && achievement.MahasiswaID != userID.(string) {
		c.JSON(403, gin.H{
			"success": false,
			"message": "Access denied: You can only upload files to your own achievements",
		})
		return
	}

	// Parse multipart form
	err = c.Request.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to parse multipart form",
			"error":   err.Error(),
		})
		return
	}

	// Get file from form
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "File is required",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	// Upload file using service
	uploadReq := service.UploadFileRequest{
		File:          file,
		FileHeader:    fileHeader,
		AchievementID: achievementID,
		UploadedBy:    userID.(string),
	}

	fileData, err := h.FileService.UploadFile(uploadReq)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Failed to upload file",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"message": "File uploaded successfully",
		"data":    fileData,
	})
}

// GetFiles godoc
func (h *AchievementHelper) GetFiles(c *gin.Context) {
	achievementID := c.Param("id")
	if achievementID == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Achievement ID is required",
		})
		return
	}

	// Check if user has access to this achievement
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	// Verify achievement exists
	achievement, err := h.AchievementService.GetAchievementByID(achievementID)
	if err != nil {
		c.JSON(404, gin.H{
			"success": false,
			"message": "Achievement not found",
			"error":   err.Error(),
		})
		return
	}

	// Check access rights
	if userRole.(string) == "mahasiswa" && achievement.MahasiswaID != userID.(string) {
		c.JSON(403, gin.H{
			"success": false,
			"message": "Access denied: You can only view files from your own achievements",
		})
		return
	}

	files, err := h.FileService.GetFiles(achievementID)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to get files",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Files retrieved successfully",
		"data": gin.H{
			"files": files,
			"total": len(files),
		},
	})
}

// DeleteFile godoc
func (h *AchievementHelper) DeleteFile(c *gin.Context) {
	achievementID := c.Param("id")
	fileID := c.Param("fileId")

	if achievementID == "" || fileID == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Achievement ID and File ID are required",
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{
			"success": false,
			"message": "User ID not found",
		})
		return
	}

	// Validate file access before deletion
	userRole, _ := c.Get("user_role")
	err := h.FileService.ValidateFileAccess(fileID, userID.(string), userRole.(string))
	if err != nil {
		c.JSON(403, gin.H{
			"success": false,
			"message": "Access denied",
			"error":   err.Error(),
		})
		return
	}

	// Delete file
	err = h.FileService.DeleteFile(fileID, userID.(string))
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to delete file",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "File deleted successfully",
	})
}

// DownloadFile godoc
func (h *AchievementHelper) DownloadFile(c *gin.Context) {
	achievementID := c.Param("id")
	fileID := c.Param("fileId")

	if achievementID == "" || fileID == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Achievement ID and File ID are required",
		})
		return
	}

	// Validate file access
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")
	err := h.FileService.ValidateFileAccess(fileID, userID.(string), userRole.(string))
	if err != nil {
		c.JSON(403, gin.H{
			"success": false,
			"message": "Access denied",
			"error":   err.Error(),
		})
		return
	}

	// Download file
	fileStream, fileData, err := h.FileService.DownloadFile(fileID)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to download file",
			"error":   err.Error(),
		})
		return
	}

	// Set appropriate headers for file download
	c.Header("Content-Type", fileData.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileData.Filename))
	c.Header("Content-Length", fmt.Sprintf("%d", fileData.Size))

	// Stream file content
	_, err = io.Copy(c.Writer, fileStream)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Failed to stream file",
			"error":   err.Error(),
		})
		return
	}
}

// SubmitAchievement godoc
func (h *AchievementHelper) SubmitAchievement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Achievement ID is required"})
		return
	}

	// Check if user is mahasiswa
	userRole, _ := c.Get("user_role")
	if userRole.(string) != "mahasiswa" {
		c.JSON(403, gin.H{"error": "Only mahasiswa can submit achievements"})
		return
	}

	err := h.AchievementService.SubmitAchievement(id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Achievement submitted for verification",
		"status":  "success",
	})
}

// SubmitAchievement godoc
func (h *AchievementHelper) VerifyAchievement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Achievement ID is required"})
		return
	}

	// Check if user is dosen or admin
	userRole, _ := c.Get("user_role")
	if userRole.(string) == "mahasiswa" {
		c.JSON(403, gin.H{"error": "Only dosen and admin can verify achievements"})
		return
	}

	userID, _ := c.Get("user_id")
	err := h.AchievementService.VerifyAchievement(id, userID.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Achievement verified successfully",
		"status":  "success",
	})
}

// RejectAchievement godoc
func (h *AchievementHelper) RejectAchievement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, gin.H{"error": "Achievement ID is required"})
		return
	}

	// Check if user is dosen or admin
	userRole, _ := c.Get("user_role")
	if userRole.(string) == "mahasiswa" {
		c.JSON(403, gin.H{"error": "Only dosen and admin can reject achievements"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Rejection reason is required"})
		return
	}

	userID, _ := c.Get("user_id")
	err := h.AchievementService.RejectAchievement(id, userID.(string), req.Reason)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Achievement rejected",
		"reason":  req.Reason,
		"status":  "success",
	})
}

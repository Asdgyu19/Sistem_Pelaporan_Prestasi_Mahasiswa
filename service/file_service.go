package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"prestasi-mahasiswa/database"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileService struct {
	MongoDB *database.MongoDB
}

type FileData struct {
	ID            string    `json:"id" bson:"_id"`
	Filename      string    `json:"filename" bson:"filename"`
	Size          int64     `json:"size" bson:"size"`
	ContentType   string    `json:"content_type" bson:"content_type"`
	UploadedAt    time.Time `json:"uploaded_at" bson:"uploaded_at"`
	AchievementID string    `json:"achievement_id" bson:"achievement_id"`
	UploadedBy    string    `json:"uploaded_by" bson:"uploaded_by"`
	GridFSID      string    `json:"gridfs_id" bson:"gridfs_id"`
}

type UploadFileRequest struct {
	File          multipart.File        `json:"-"`
	FileHeader    *multipart.FileHeader `json:"-"`
	AchievementID string                `json:"achievement_id"`
	UploadedBy    string                `json:"uploaded_by"`
}

func NewFileService(mongodb *database.MongoDB) *FileService {
	return &FileService{
		MongoDB: mongodb,
	}
}

// UploadFile uploads file to MongoDB GridFS with metadata
func (s *FileService) UploadFile(req UploadFileRequest) (*FileData, error) {
	if req.File == nil || req.FileHeader == nil {
		return nil, errors.New("file is required")
	}

	if req.AchievementID == "" {
		return nil, errors.New("achievement_id is required")
	}

	if req.UploadedBy == "" {
		return nil, errors.New("uploaded_by is required")
	}

	// Validate file type
	if !s.isValidFileType(req.FileHeader.Filename) {
		return nil, errors.New("invalid file type. Allowed: pdf, doc, docx, jpg, jpeg, png")
	}

	// Validate file size (5MB limit)
	if req.FileHeader.Size > 5*1024*1024 {
		return nil, errors.New("file size exceeds 5MB limit")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize GridFS bucket
	bucket, err := gridfs.NewBucket(s.MongoDB.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to create GridFS bucket: %v", err)
	}

	// Generate unique filename
	fileID := uuid.New().String()
	uniqueFilename := fmt.Sprintf("%s_%s", fileID, req.FileHeader.Filename)

	// Upload file to GridFS
	uploadStream, err := bucket.OpenUploadStreamWithID(
		primitive.NewObjectID(),
		uniqueFilename,
		options.GridFSUpload().SetMetadata(bson.D{
			{"achievement_id", req.AchievementID},
			{"uploaded_by", req.UploadedBy},
			{"original_filename", req.FileHeader.Filename},
			{"content_type", s.getContentType(req.FileHeader.Filename)},
			{"uploaded_at", time.Now()},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload stream: %v", err)
	}
	defer uploadStream.Close()

	// Copy file content to GridFS
	_, err = io.Copy(uploadStream, req.File)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	// Create file metadata document
	fileData := &FileData{
		ID:            fileID,
		Filename:      req.FileHeader.Filename,
		Size:          req.FileHeader.Size,
		ContentType:   s.getContentType(req.FileHeader.Filename),
		UploadedAt:    time.Now(),
		AchievementID: req.AchievementID,
		UploadedBy:    req.UploadedBy,
		GridFSID:      uploadStream.FileID.(primitive.ObjectID).Hex(),
	}

	// Store file metadata in collection
	collection := s.MongoDB.Database.Collection("achievement_files")
	_, err = collection.InsertOne(ctx, fileData)
	if err != nil {
		return nil, fmt.Errorf("failed to store file metadata: %v", err)
	}

	return fileData, nil
}

// GetFiles retrieves all files for an achievement
func (s *FileService) GetFiles(achievementID string) ([]FileData, error) {
	if achievementID == "" {
		return nil, errors.New("achievement_id is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.MongoDB.Database.Collection("achievement_files")
	filter := bson.M{"achievement_id": achievementID}

	cursor, err := collection.Find(ctx, filter, options.Find().SetSort(bson.D{{"uploaded_at", -1}}))
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %v", err)
	}
	defer cursor.Close(ctx)

	var files []FileData
	if err = cursor.All(ctx, &files); err != nil {
		return nil, fmt.Errorf("failed to decode files: %v", err)
	}

	return files, nil
}

// GetFileByID retrieves specific file metadata
func (s *FileService) GetFileByID(fileID string) (*FileData, error) {
	if fileID == "" {
		return nil, errors.New("file_id is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.MongoDB.Database.Collection("achievement_files")
	filter := bson.M{"_id": fileID}

	var fileData FileData
	err := collection.FindOne(ctx, filter).Decode(&fileData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("file not found")
		}
		return nil, fmt.Errorf("failed to get file: %v", err)
	}

	return &fileData, nil
}

// DownloadFile streams file from GridFS
func (s *FileService) DownloadFile(fileID string) (io.Reader, *FileData, error) {
	// Get file metadata
	fileData, err := s.GetFileByID(fileID)
	if err != nil {
		return nil, nil, err
	}

	// Initialize GridFS bucket
	bucket, err := gridfs.NewBucket(s.MongoDB.Database)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create GridFS bucket: %v", err)
	}

	// Convert hex string back to ObjectID
	gridfsID, err := primitive.ObjectIDFromHex(fileData.GridFSID)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid GridFS ID: %v", err)
	}

	// Open download stream
	downloadStream, err := bucket.OpenDownloadStream(gridfsID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open download stream: %v", err)
	}

	return downloadStream, fileData, nil
}

// DeleteFile removes file from both GridFS and metadata collection
func (s *FileService) DeleteFile(fileID, userID string) error {
	if fileID == "" {
		return errors.New("file_id is required")
	}

	if userID == "" {
		return errors.New("user_id is required")
	}

	// Get file metadata first
	fileData, err := s.GetFileByID(fileID)
	if err != nil {
		return err
	}

	// Check if user has permission to delete (owner or admin)
	if fileData.UploadedBy != userID {
		// TODO: Add admin role check here
		return errors.New("permission denied: you can only delete your own files")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete from GridFS
	bucket, err := gridfs.NewBucket(s.MongoDB.Database)
	if err != nil {
		return fmt.Errorf("failed to create GridFS bucket: %v", err)
	}

	gridfsID, err := primitive.ObjectIDFromHex(fileData.GridFSID)
	if err != nil {
		return fmt.Errorf("invalid GridFS ID: %v", err)
	}

	err = bucket.Delete(gridfsID)
	if err != nil {
		return fmt.Errorf("failed to delete file from GridFS: %v", err)
	}

	// Delete metadata
	collection := s.MongoDB.Database.Collection("achievement_files")
	_, err = collection.DeleteOne(ctx, bson.M{"_id": fileID})
	if err != nil {
		return fmt.Errorf("failed to delete file metadata: %v", err)
	}

	return nil
}

// ValidateFileAccess checks if user can access file
func (s *FileService) ValidateFileAccess(fileID, userID, userRole string) error {
	fileData, err := s.GetFileByID(fileID)
	if err != nil {
		return err
	}

	// Owner can always access
	if fileData.UploadedBy == userID {
		return nil
	}

	// Admin and dosen can access all files
	if userRole == "admin" || userRole == "dosen_wali" {
		return nil
	}

	return errors.New("access denied")
}

// Helper functions
func (s *FileService) isValidFileType(filename string) bool {
	allowedExtensions := []string{".pdf", ".doc", ".docx", ".jpg", ".jpeg", ".png"}
	ext := strings.ToLower(filepath.Ext(filename))

	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

func (s *FileService) getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
}

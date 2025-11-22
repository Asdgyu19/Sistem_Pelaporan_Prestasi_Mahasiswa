package service

import (
	"prestasi-mahasiswa/database"
)

type FileService struct {
	MongoDB *database.MongoDB
}

func NewFileService(mongodb *database.MongoDB) *FileService {
	return &FileService{
		MongoDB: mongodb,
	}
}

func (s *FileService) UploadFile(achievementID string, filename string, data []byte) error {
	// TODO: Implement file upload to MongoDB
	return nil
}

func (s *FileService) GetFiles(achievementID string) ([]FileData, error) {
	// TODO: Implement get files from MongoDB
	return []FileData{}, nil
}

func (s *FileService) DeleteFile(achievementID, fileID string) error {
	// TODO: Implement file deletion from MongoDB
	return nil
}

type FileData struct {
	ID            string `json:"id"`
	Filename      string `json:"filename"`
	Size          int64  `json:"size"`
	ContentType   string `json:"content_type"`
	UploadedAt    string `json:"uploaded_at"`
	AchievementID string `json:"achievement_id"`
}

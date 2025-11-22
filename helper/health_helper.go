package helper

import (
	"context"
	"database/sql"
	"prestasi-mahasiswa/database"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHelper struct {
	DB      *sql.DB
	MongoDB *database.MongoDB
}

func NewHealthHelper(db *sql.DB, mongodb *database.MongoDB) *HealthHelper {
	return &HealthHelper{
		DB:      db,
		MongoDB: mongodb,
	}
}

func (h *HealthHelper) CheckHealth(c *gin.Context) {
	// Test PostgreSQL connection
	pgStatus := "connected"
	if err := h.DB.Ping(); err != nil {
		pgStatus = "disconnected: " + err.Error()
	}

	// Test MongoDB connection
	mongoStatus := "connected"
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := h.MongoDB.Client.Ping(ctx, nil); err != nil {
		mongoStatus = "disconnected: " + err.Error()
	}

	c.JSON(200, gin.H{
		"status":    "OK",
		"timestamp": time.Now().Format(time.RFC3339),
		"database": gin.H{
			"postgresql": pgStatus,
			"mongodb":    mongoStatus,
		},
		"services": gin.H{
			"api":  "running",
			"port": "8080",
		},
	})
}

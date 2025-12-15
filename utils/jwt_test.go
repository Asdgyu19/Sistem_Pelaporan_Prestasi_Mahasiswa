package utils

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	jwtUtil := NewJWTUtil("test-secret-key-that-is-long-enough-for-testing", 24)

	userID := "test-user-123"
	email := "test@example.com"
	role := "mahasiswa"

	token, err := jwtUtil.GenerateToken(userID, email, role)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Expected token, got empty string")
	}
}

func TestValidateToken(t *testing.T) {
	secretKey := "test-secret-key-that-is-long-enough-for-testing"
	jwtUtil := NewJWTUtil(secretKey, 24)

	userID := "test-user-123"
	email := "test@example.com"
	role := "mahasiswa"

	token, _ := jwtUtil.GenerateToken(userID, email, role)

	claims, err := jwtUtil.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected userID %s, got %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Expected email %s, got %s", email, claims.Email)
	}

	if claims.Role != role {
		t.Errorf("Expected role %s, got %s", role, claims.Role)
	}
}

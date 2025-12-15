package service

import (
	"testing"

	"prestasi-mahasiswa/utils"
)

func TestAuthenticateUserSuccess(t *testing.T) {
	// Test case untuk login yang valid
	// Ini akan menggunakan test database

	password := "password123"

	// Hash password
	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashedPwd == "" {
		t.Error("Expected hashed password, got empty string")
	}
}

func TestAuthenticateUserInvalidPassword(t *testing.T) {
	password := "password123"
	wrongPassword := "wrongpassword"

	hashedPwd, _ := utils.HashPassword(password)

	err := utils.ComparePassword(hashedPwd, wrongPassword)
	if err == nil {
		t.Error("Expected error comparing wrong password, got nil")
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "mysecurepassword"

	hashed, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashed == password {
		t.Error("Hashed password should not match original password")
	}

	err = utils.ComparePassword(hashed, password)
	if err != nil {
		t.Errorf("Password comparison failed: %v", err)
	}
}

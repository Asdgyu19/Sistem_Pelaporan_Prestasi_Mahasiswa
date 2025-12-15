package utils

import (
	"testing"
)

func TestHashPasswordSuccess(t *testing.T) {
	password := "securepassword123"

	hashedPwd, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hashedPwd == "" {
		t.Error("Expected hashed password, got empty string")
	}

	if hashedPwd == password {
		t.Error("Hashed password should not equal original password")
	}
}

func TestComparePasswordSuccess(t *testing.T) {
	password := "securepassword123"

	hashedPwd, _ := HashPassword(password)

	err := ComparePassword(hashedPwd, password)
	if err != nil {
		t.Fatalf("Failed to compare password: %v", err)
	}
}

func TestComparePasswordFailed(t *testing.T) {
	password := "securepassword123"
	wrongPassword := "wrongpassword"

	hashedPwd, _ := HashPassword(password)

	err := ComparePassword(hashedPwd, wrongPassword)
	if err == nil {
		t.Error("Expected error when comparing wrong password, got nil")
	}
}

func TestHashPasswordDifferent(t *testing.T) {
	password := "securepassword123"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// Setiap hash harus unik (karena salt yang berbeda)
	if hash1 == hash2 {
		t.Error("Two hashes of the same password should be different")
	}

	// Tapi keduanya harus cocok saat di-compare
	if err := ComparePassword(hash1, password); err != nil {
		t.Errorf("Hash1 comparison failed: %v", err)
	}

	if err := ComparePassword(hash2, password); err != nil {
		t.Errorf("Hash2 comparison failed: %v", err)
	}
}

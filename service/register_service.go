package service

import (
	"database/sql"
	"errors"
)

type RegisterService struct {
	DB *sql.DB
}

func NewRegisterService(db *sql.DB) *RegisterService {
	return &RegisterService{DB: db}
}

func (s *RegisterService) CreateUser(nim, name, email, password, role string) error {
	// TODO: Implement user creation logic
	if email == "" || password == "" || name == "" {
		return errors.New("all fields are required")
	}

	// TODO: Hash password, insert into database
	// For now just placeholder
	return nil
}

func (s *RegisterService) ValidateUserData(nim, name, email, password, role string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	// TODO: Add more validation logic
	return nil
}

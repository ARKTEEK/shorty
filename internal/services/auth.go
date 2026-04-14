package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ARKTEEK/shorty/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sql.DB
	us *UserService
}

func NewAuthService(db *sql.DB, us *UserService) *AuthService {
	return &AuthService{db: db, us: us}
}

func (s *AuthService) Login(ctx context.Context, request models.AuthRequest) (*models.LoginResponse, error) {
	exists, err := s.us.Exists(ctx, request.Email)
	if err != nil {
		return nil, fmt.Errorf("check existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("invalid credentials")
	}

	var storedHash string
	err = s.db.QueryRowContext(ctx, "SELECT password FROM users WHERE email = ?", request.Email).Scan(&storedHash)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(request.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	return &models.LoginResponse{
		UserID:  1,
		Email:   request.Email,
		Token:   fmt.Sprintf("token-%d", 12345),
		Message: "Login successful",
	}, nil

}

func (s *AuthService) Register(ctx context.Context, request models.AuthRequest) (*models.RegisterResponse, error) {
	exists, err := s.us.Exists(ctx, request.Email)
	if err != nil {
		return nil, fmt.Errorf("check existence: %w", err)
	}

	if exists {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	result, err := s.db.ExecContext(ctx,
		`INSERT INTO users (email, password) VALUES (?, ?)`,
		request.Email, string(hash),
	)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	id, _ := result.LastInsertId()
	return &models.RegisterResponse{
		ID:       id,
		Username: request.Email,
		Message:  "Registration successful",
	}, nil
}

func (s *AuthService) Deactivate(ctx context.Context, request models.DeactivateRequest) (*models.DeactivateResponse, error) {
	exists, err := s.us.Exists(ctx, request.Email)
	if err != nil {
		return nil, fmt.Errorf("check existence: %w", err)
	}

	if !exists {
		return nil, errors.New("User not found!")
	}

	_, err = s.db.ExecContext(ctx,
		`UPDATE users SET active = false WHERE email = ?`,
		request.Email)

	if err != nil {
		return nil, fmt.Errorf("deactivate user: %w", err)
	}

	return &models.DeactivateResponse{
		Success: true,
		Message: "User deactivated successfully",
	}, nil
}

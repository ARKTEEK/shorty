package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ARKTEEK/shorty/internal/middleware"
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
		return nil, fmt.Errorf("Check existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("Invalid credentials.")
	}

	var (
		userID     int64
		storedHash string
	)
	err = s.db.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = ? LIMIT 1", request.Email).
		Scan(&userID, &storedHash)
	if err != nil {
		return nil, fmt.Errorf("Invalid credentials.")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(request.Password)); err != nil {
		return nil, fmt.Errorf("Invalid credentials.")
	}

	token, err := middleware.GenerateToken(userID)
	if err != nil {
		return nil, fmt.Errorf("Generate token: %w", err)
	}

	return &models.LoginResponse{
		UserID:  userID,
		Email:   request.Email,
		Token:   token,
		Message: "Login successful.",
	}, nil
}

func (s *AuthService) Register(ctx context.Context, request models.AuthRequest) (*models.RegisterResponse, error) {
	exists, err := s.us.Exists(ctx, request.Email)
	if err != nil {
		return nil, fmt.Errorf("Check existence: %w", err)
	}

	if exists {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Hash password: %w", err)
	}

	result, err := s.db.ExecContext(ctx,
		`INSERT INTO users (email, password) VALUES (?, ?)`,
		request.Email, string(hash),
	)
	if err != nil {
		return nil, fmt.Errorf("Insert user: %w", err)
	}

	id, _ := result.LastInsertId()
	return &models.RegisterResponse{
		ID:       id,
		Username: request.Email,
		Message:  "Registration successful.",
	}, nil
}

func (s *AuthService) Deactivate(ctx context.Context, request *models.DeactivateRequest) (*models.DeactivateResponse, error) {
	var storedHash string

	err := s.db.QueryRowContext(ctx,
		`SELECT password FROM users WHERE id = ? AND active = true LIMIT 1`,
		request.UserID,
	).Scan(&storedHash)

	if err != nil {
		return nil, errors.New("Invalid credentials.")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(request.Password)); err != nil {
		return nil, errors.New("Invalid credentials.")
	}

	_, err = s.db.ExecContext(ctx,
		`UPDATE users SET active = false WHERE id = ?`,
		request.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("Deactivate user: %w", err)
	}

	return &models.DeactivateResponse{
		Success: true,
		Message: "User deactivated successfully",
	}, nil
}

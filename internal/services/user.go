package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ARKTEEK/shorty/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailTaken = errors.New("email already in use")

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetById(ctx context.Context, id int64) (*models.User, error) {
	u := &models.User{}
	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, created_at FROM users WHERE id = ?`,
		id,
	).Scan(&u.ID, &u.Email, &u.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	return u, nil
}

func (s *UserService) Exists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`
	err := s.db.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}

func (s *UserService) Update(ctx context.Context, id int64, request models.UpdateUserRequest) (*models.User, error) {
	existingUser, err := s.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if request.Email != "" && request.Email != existingUser.Email {
		exists, err := s.Exists(ctx, request.Email)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already taken")
		}
	}

	query := "UPDATE users SET "
	args := []any{}

	if request.Email != "" {
		query += "email = ?, "
		args = append(args, request.Email)
	}

	if request.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		query += "password = ?, "
		args = append(args, hashedPassword)
	}

	query = query[:len(query)-2] + " WHERE id = ?"
	args = append(args, id)

	_, err = s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("db update failed: %w", err)
	}

	return s.GetById(ctx, id)
}

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

func (s *UserService) Create(ctx context.Context, request models.CreateUserRequest) (*models.User, error) {
	exists, err := s.Exists(ctx, request.Email)
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
	return s.GetById(ctx, id)
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

func (s *UserService) Deactivate(ctx context.Context, request models.DeactivateUserRequest) (bool, error) {
	exists, err := s.Exists(ctx, request.Email)
	if err != nil {
		return false, fmt.Errorf("check existence: %w", err)
	}

	if !exists {
		return false, errors.New("User not found!")
	}

	_, err = s.db.ExecContext(ctx,
		`UPDATE users SET active = false WHERE email = ?`,
		request.Email)

	if err != nil {
		return false, fmt.Errorf("deactivate user: %w", err)
	}

	return true, err
}

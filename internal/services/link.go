package services

import (
	"context"
	"database/sql"
	"math/rand"

	"github.com/ARKTEEK/shorty/internal/models"
)

type LinkService struct {
	db *sql.DB
	us *UserService
}

func NewLinkService(db *sql.DB, us *UserService) *LinkService {
	return &LinkService{db: db, us: us}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString() string {
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (s *LinkService) Exists(ctx context.Context, shortCode string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM links WHERE short_code = ?)`
	err := s.db.QueryRowContext(ctx, query, shortCode).Scan(&exists)
	return exists, err
}

func (s *LinkService) CreateShortLink(ctx context.Context, req *models.CreateLinkRequest) (*models.Link, error) {
	shortCode := generateRandomString()

	if _, err := s.Exists(ctx, shortCode); err != nil {
		return nil, err
	}

	_, err := s.db.Exec("INSERT INTO links (original_url, short_code, user_id) VALUES (?, ?, ?)", req.OriginalUrl, shortCode, req.UserId)
	if err != nil {
		return nil, err
	}

	return &models.Link{
		OriginalUrl: req.OriginalUrl,
		ShortCode:   shortCode,
		UserID:      req.UserId,
		Visits:      0,
	}, nil
}

func (s *LinkService) GetOriginalUrl(ctx context.Context, shortCode string) (string, error) {
	var originalUrl string
	err := s.db.QueryRowContext(ctx, "SELECT original_url FROM links WHERE short_code = ?", shortCode).Scan(&originalUrl)
	return originalUrl, err
}

func (s *LinkService) IncrementVisits(ctx context.Context, shortCode string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE links SET visits = visits + 1 WHERE short_code = ?", shortCode)
	return err
}

func (s *LinkService) ListLinksByUser(ctx context.Context, userID int32) ([]models.MyLink, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, original_url, short_code, visits, created_at
		 FROM links
		 WHERE user_id = ?
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.MyLink
	for rows.Next() {
		var l models.MyLink
		if err := rows.Scan(&l.ID, &l.OriginalUrl, &l.ShortCode, &l.Visits, &l.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *LinkService) DeleteLink(ctx context.Context, userID int32, shortCode string) (bool, error) {
	res, err := s.db.ExecContext(
		ctx,
		`DELETE FROM links WHERE short_code = ? AND user_id = ?`,
		shortCode,
		userID,
	)
	if err != nil {
		return false, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

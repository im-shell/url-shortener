package services

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

type urlService struct {
	DB *sql.DB
}

type URLService interface {
	IsActive(key string) (bool, error)
	ShortenURL(longURL string) (string, error)
	RedirectURL(key string) (string, error)
}

func NewURLService(db *sql.DB) *urlService {
	return &urlService{DB: db}
}

func (s *urlService) IsActive(key string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM urls WHERE key = $1 AND expires_at > NOW())"
	var exists bool
	err := s.DB.QueryRow(query, key).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *urlService) ShortenURL(longURL string) (string, error) {
	// 1. Generate a unique key
	// 2. Check if the key is already in the database
	// 3. If the key is not in the database, insert the long URL and key into the database
	// 4. Return the short URL

	query := "SELECT EXISTS(SELECT 1 from urls WHERE key = $1)"

	var existingKey string
	var key string

	for {
		key = uuid.New().String()[:6]
		err := s.DB.QueryRow(query, key).Scan(&existingKey)
		if err != nil {
			return "", err
		}
		if existingKey == "t" {
			continue
		}
		query = "INSERT INTO urls (long_url, key, created_at, expires_at) VALUES ($1, $2, $3, $4)"
		_, err = s.DB.Exec(query, longURL, key, time.Now(), time.Now().Add(24*time.Hour))
		if err != nil {
			return "", err
		}
		break
	}

	shortURL := fmt.Sprintf("%s/%s", os.Getenv("BASE_URL"), key)

	return shortURL, nil
}

func (s *urlService) RedirectURL(key string) (string, error) {
	query := "SELECT long_url FROM urls WHERE key = $1"
	var longURL string
	err := s.DB.QueryRow(query, key).Scan(&longURL)
	if err != nil {
		return "", err
	}

	_, err = s.DB.Exec(query, longURL, key, time.Now(), time.Now().Add(24*time.Hour))
	if err != nil {
		return "", err
	}

	return key, nil
}

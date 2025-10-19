package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UrlRepository struct {
	pgxPool *pgxpool.Pool
}

func NewUrlRepository(pgxPool *pgxpool.Pool) *UrlRepository {
	return &UrlRepository{pgxPool: pgxPool}
}

func (repo *UrlRepository) Initialize() error {
	createTable := `
	CREATE TABLE IF NOT EXISTS urls (
		short_url VARCHAR(16) PRIMARY KEY,
		long_url VARCHAR(2048) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)
	`
	log.Println("Initializing database")

	_, err := repo.pgxPool.Exec(context.Background(), createTable)
	return err
}

func (repo *UrlRepository) Save(mapping UrlMap) error {
	insertRow := `
	INSERT INTO urls (short_url, long_url) VALUES ($1, $2)
	`
	log.Println("Saving URL to DB")

	_, err := repo.pgxPool.Exec(context.Background(), insertRow, mapping.ShortUrl, mapping.LongUrl)
	return err
}

func (repo *UrlRepository) Get(shortUrl string) (string, bool, error) {
	selectRow := `
	SELECT long_url FROM urls WHERE short_url = $1 LIMIT 1
	`
	log.Println("Finding URL from DB")

	row := repo.pgxPool.QueryRow(context.Background(), selectRow, shortUrl)

	var longUrl string
	err := row.Scan(&longUrl)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}

		return "", false, err
	}

	return longUrl, true, nil
}

package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"

	"scan-service/models"
)

const (
	InsertScan        = "INSERT INTO scan_requests(url, is_private, encrypted_token) values ($1, $2, $3) returning id"
	UpdateQueueStatus = "UPDATE scan_requests SET queue_status = $1 where id = $2"
)

type ScanRepository interface {
	InsertScan(ctx context.Context, scanRequest models.ScanRequest) (*string, error)
	UpdateQueueStatus(ctx context.Context, scanID string, status models.QueueStatus) error
}

type scanRepository struct {
	db *sqlx.DB
}

func (s scanRepository) InsertScan(ctx context.Context, scanRequest models.ScanRequest) (*string, error) {
	log.Println("Inserting scan to scan request table")
	var requestID string
	err := s.db.QueryRowContext(ctx, InsertScan, scanRequest.URL, scanRequest.IsPrivate, scanRequest.EncryptedToken).Scan(&requestID)
	if err != nil {
		log.Printf("Unable to insert scan request %v", err)
		return nil, err
	}

	log.Println("returning scan request id ", requestID)

	return &requestID, nil
}

func (s scanRepository) UpdateQueueStatus(ctx context.Context, scanID string, status models.QueueStatus) error {
	log.Println("Updating scan queue status")
	_, err := s.db.ExecContext(ctx, UpdateQueueStatus, status, scanID)
	if err != nil {
		log.Printf("Unable to insert scan request %v", err)
		return err
	}

	log.Println("successfully updated the queue status for scan id ", scanID)

	return nil
}

func NewScanRepository(db *sqlx.DB) ScanRepository {
	return &scanRepository{db: db}
}

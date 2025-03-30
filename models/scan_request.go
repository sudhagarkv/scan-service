package models

import (
	"time"

	"scan-service/scm"
)

type QueueStatus string

const (
	NotQueued QueueStatus = "NotQueued"
	Queued    QueueStatus = "Queued"
	Scheduled QueueStatus = "Scheduled"
	Processed QueueStatus = "Processed"
)

type ScanRequest struct {
	ID             string      `json:"id,omitempty" db:"id"`
	URL            string      `json:"url" binding:"required" db:"url"`
	Type           scm.Type    `json:"type"`
	IsPrivate      bool        `json:"isPrivate" db:"is_private"`
	EncryptedToken string      `json:"encryptedToken,omitempty" binding:"isPrivate" db:"encrypted_token"`
	ScanID         string      `json:"scanID" db:"scan_id"`
	QueueStatus    QueueStatus `json:"queueStatus" db:"queue_status"`
	CreatedAt      time.Time   `json:"createdAt" db:"created_at"`
	ModifiedAt     time.Time   `json:"modifiedAt" db:"modified_at"`
	DeletedAt      time.Time   `json:"deletedAt" db:"deleted_at"`
}

type QueuedScan struct {
	RepoName       string
	IsPrivate      bool
	EncryptedToken string
	Namespace      string
	URL            string
	ID             string
}

package views

import (
	"github.com/google/uuid"
)

type FormattedLog struct {
	ID        string
	Action    string
	AccountID string
	Details   map[string]interface{}
	CreatedAt string
}

type Pagination struct {
	CurrentPage uint64
	TotalPages  uint64
	Limit       uint64
	Offset      uint64
	HasMore     bool
	TotalItems  uint64
}

type ActivityLogs struct {
	AccountID   uuid.UUID
	AccountName string
	Logs        []FormattedLog
	Pagination  Pagination
}

package dto

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/build/domain"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// ListBuildEventFilters defines filters for listing build events
type ListBuildEventFilters struct {
	ProjectID *value_objects.ID
	EventType *domain.EventType
	Status    *domain.BuildStatus
	Branch    *string
	DateFrom  *time.Time
	DateTo    *time.Time
	Limit     int
	Offset    int
	OrderBy   string
	OrderDir  string // "asc" or "desc"
}

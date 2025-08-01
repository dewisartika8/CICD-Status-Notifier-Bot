package postgres

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/dashboard/port"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
)

// dashboardProjectRepository implements dashboard specific project queries
type dashboardProjectRepository struct {
	db *gorm.DB
}

// NewDashboardProjectRepository creates a new dashboard project repository
func NewDashboardProjectRepository(db *gorm.DB) port.ProjectRepositoryInterface {
	return &dashboardProjectRepository{db: db}
}

// GetProjectNameByID returns project name by ID
func (r *dashboardProjectRepository) GetProjectNameByID(projectID value_objects.ID) (string, error) {
	var name string
	query := "SELECT name FROM projects WHERE id = ? AND is_active = true"

	if err := r.db.Raw(query, projectID.Value()).Scan(&name).Error; err != nil {
		return "", fmt.Errorf("failed to get project name: %w", err)
	}

	if name == "" {
		return "", fmt.Errorf("project not found")
	}

	return name, nil
}

// GetActiveProjectsCount returns count of active projects
func (r *dashboardProjectRepository) GetActiveProjectsCount() (int, error) {
	var count int64
	if err := r.db.Model(&struct {
		ID       string `gorm:"column:id;primaryKey"`
		IsActive bool   `gorm:"column:is_active"`
	}{}).Where("is_active = ?", true).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get active projects count: %w", err)
	}

	return int(count), nil
}

// GetTotalProjectsCount returns total count of projects
func (r *dashboardProjectRepository) GetTotalProjectsCount() (int, error) {
	var count int64
	if err := r.db.Model(&struct {
		ID string `gorm:"column:id;primaryKey"`
	}{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to get total projects count: %w", err)
	}

	return int(count), nil
}

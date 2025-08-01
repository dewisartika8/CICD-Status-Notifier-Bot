package dto

import (
	"time"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/webhook/domain"
)

// GitHubActionsPayload represents GitHub webhook payload structure
type GitHubActionsPayload struct {
	Action     string `json:"action"`
	Workflow   string `json:"workflow"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		HTMLURL  string `json:"html_url"`
	} `json:"repository"`
	Sender struct {
		Login string `json:"login"`
		Email string `json:"email,omitempty"`
	} `json:"sender"`
	WorkflowRun *WorkflowRun `json:"workflow_run,omitempty"`

	// Push event specific fields
	Ref        string   `json:"ref,omitempty"`
	Before     string   `json:"before,omitempty"`
	After      string   `json:"after,omitempty"`
	Created    bool     `json:"created,omitempty"`
	Deleted    bool     `json:"deleted,omitempty"`
	Forced     bool     `json:"forced,omitempty"`
	Compare    string   `json:"compare,omitempty"`
	Commits    []Commit `json:"commits,omitempty"`
	HeadCommit *Commit  `json:"head_commit,omitempty"`
	Pusher     *User    `json:"pusher,omitempty"`

	// Pull request event specific fields
	Number      int          `json:"number,omitempty"`
	PullRequest *PullRequest `json:"pull_request,omitempty"`
}

// WorkflowRun represents GitHub workflow run information
type WorkflowRun struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	Conclusion string    `json:"conclusion"`
	HTMLURL    string    `json:"html_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	RunNumber  int       `json:"run_number"`
	Event      string    `json:"event"`
	HeadBranch string    `json:"head_branch"`
	HeadSha    string    `json:"head_sha"`
}

// User represents a GitHub user
type User struct {
	Login string `json:"login"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

// Commit represents a git commit
type Commit struct {
	ID        string   `json:"id"`
	TreeID    string   `json:"tree_id"`
	Distinct  bool     `json:"distinct"`
	Message   string   `json:"message"`
	Timestamp string   `json:"timestamp"`
	URL       string   `json:"url"`
	Author    User     `json:"author"`
	Committer User     `json:"committer"`
	Added     []string `json:"added,omitempty"`
	Removed   []string `json:"removed,omitempty"`
	Modified  []string `json:"modified,omitempty"`
}

// PullRequest represents a GitHub pull request
type PullRequest struct {
	ID               int        `json:"id"`
	Number           int        `json:"number"`
	State            string     `json:"state"`
	Title            string     `json:"title"`
	Body             string     `json:"body"`
	HTMLURL          string     `json:"html_url"`
	DiffURL          string     `json:"diff_url"`
	PatchURL         string     `json:"patch_url"`
	IssueURL         string     `json:"issue_url"`
	CommitsURL       string     `json:"commits_url"`
	ReviewCommentURL string     `json:"review_comment_url"`
	StatusesURL      string     `json:"statuses_url"`
	User             User       `json:"user"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	ClosedAt         *time.Time `json:"closed_at,omitempty"`
	MergedAt         *time.Time `json:"merged_at,omitempty"`
	Assignee         *User      `json:"assignee,omitempty"`
	Assignees        []User     `json:"assignees,omitempty"`
	Merged           bool       `json:"merged"`
	Mergeable        *bool      `json:"mergeable,omitempty"`
	Head             Branch     `json:"head"`
	Base             Branch     `json:"base"`
}

// Branch represents a git branch
type Branch struct {
	Label string     `json:"label"`
	Ref   string     `json:"ref"`
	SHA   string     `json:"sha"`
	User  User       `json:"user"`
	Repo  Repository `json:"repo"`
}

// Repository represents a GitHub repository
type Repository struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    User   `json:"owner"`
	Private  bool   `json:"private"`
	HTMLURL  string `json:"html_url"`
	CloneURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
}

// ProcessWebhookRequest represents a request to process a webhook
type ProcessWebhookRequest struct {
	ProjectID  value_objects.ID        `json:"project_id" validate:"required"`
	EventType  domain.WebhookEventType `json:"event_type" validate:"required"`
	Signature  string                  `json:"signature" validate:"required"`
	DeliveryID string                  `json:"delivery_id"`
	Body       []byte                  `json:"body" validate:"required"`
	Payload    GitHubActionsPayload    `json:"payload"`
}

// WebhookEventResponse represents webhook event response
type WebhookEventResponse struct {
	ID          string                  `json:"id"`
	ProjectID   string                  `json:"project_id"`
	EventType   domain.WebhookEventType `json:"event_type"`
	DeliveryID  string                  `json:"delivery_id"`
	ProcessedAt *time.Time              `json:"processed_at"`
	CreatedAt   time.Time               `json:"created_at"`
}

// ToWebhookEventResponse converts domain entity to response DTO
func ToWebhookEventResponse(event *domain.WebhookEvent) *WebhookEventResponse {
	return &WebhookEventResponse{
		ID:          event.ID().String(),
		ProjectID:   event.ProjectID().String(),
		EventType:   event.EventType(),
		DeliveryID:  event.DeliveryID(),
		ProcessedAt: event.ProcessedAt(),
		CreatedAt:   event.CreatedAt().ToTime(),
	}
}

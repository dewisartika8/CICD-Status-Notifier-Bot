package repository

// Query constants to avoid duplication
const (
	queryByID            = "id = ?"
	queryByName          = "name = ?"
	queryByProjectID     = "project_id = ?"
	queryByEventType     = "event_type = ?"
	queryByBuildEventID  = "build_event_id = ?"
	queryByRecipient     = "recipient = ?"
	queryByStatus        = "status = ?"
	queryByBranch        = "branch = ?"
	queryByTemplateType  = "template_type = ?"
	queryByChannel       = "channel = ?"
	queryByIsActive      = "is_active = ?"
	queryByRepositoryURL = "repository_url = ?"
	queryCreatedAtGTE    = "created_at >= ?"
	queryCreatedAtLTE    = "created_at <= ?"
	orderByCreatedAtDesc = "created_at DESC"
	orderByNameAsc       = "name ASC"
)

package dto

// StartCommandRequest represents a start command request
type StartCommandRequest struct {
	ChatID        int64  `json:"chat_id"`
	UserFirstName string `json:"user_first_name"`
}

// StartCommandResponse represents the response for /start command
type StartCommandResponse struct {
	WelcomeMessage string `json:"welcome_message"`
	UserFirstName  string `json:"user_first_name"`
	ChatID         int64  `json:"chat_id"`
}

// HelpCommandResponse represents the response for /help command
type HelpCommandResponse struct {
	HelpText      string         `json:"help_text"`
	Commands      []CommandHelp  `json:"commands"`
	UsageExamples []UsageExample `json:"usage_examples"`
}

// CommandHelp represents help information for a command
type CommandHelp struct {
	Command     string `json:"command"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
	Category    string `json:"category"`
}

// UsageExample represents a usage example
type UsageExample struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// StatusCommandRequest represents a status command request
type StatusCommandRequest struct {
	ProjectName string `json:"project_name,omitempty"`
	ChatID      int64  `json:"chat_id"`
	UserID      int64  `json:"user_id"`
}

// StatusCommandResponse represents the response for /status command
type StatusCommandResponse struct {
	ProjectName string     `json:"project_name"`
	Status      string     `json:"status"`
	Message     string     `json:"message"`
	LastBuild   *BuildInfo `json:"last_build,omitempty"`
}

// BuildInfo represents build information
type BuildInfo struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Branch    string `json:"branch"`
	CommitSHA string `json:"commit_sha"`
	Author    string `json:"author"`
	Duration  int    `json:"duration_seconds"`
	BuildURL  string `json:"build_url"`
}

// SubscribeCommandRequest represents a subscribe command request
type SubscribeCommandRequest struct {
	ProjectName string `json:"project_name"`
	ChatID      int64  `json:"chat_id"`
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
}

// SubscribeCommandResponse represents the response for /subscribe command
type SubscribeCommandResponse struct {
	ProjectName string `json:"project_name"`
	Success     bool   `json:"success"`
	Message     string `json:"message"`
}

// UnsubscribeCommandRequest represents an unsubscribe command request
type UnsubscribeCommandRequest struct {
	ProjectName string `json:"project_name"`
	ChatID      int64  `json:"chat_id"`
	UserID      int64  `json:"user_id"`
}

// UnsubscribeCommandResponse represents the response for /unsubscribe command
type UnsubscribeCommandResponse struct {
	ProjectName string `json:"project_name"`
	Success     bool   `json:"success"`
	Message     string `json:"message"`
}

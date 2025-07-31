package domain

import (
	"errors"
	"strings"
)

// CommandContext represents the context of a bot command
type CommandContext struct {
	Command  string
	Args     []string
	UserID   int64
	ChatID   int64
	Username string
}

// CommandValidator handles command validation
type CommandValidator struct {
	allowedUsers map[int64]bool
}

// NewCommandValidator creates a new command validator
func NewCommandValidator() *CommandValidator {
	return &CommandValidator{
		allowedUsers: make(map[int64]bool),
	}
}

// ValidateCommand validates the command context
func (cv *CommandValidator) ValidateCommand(ctx *CommandContext) error {
	// Validate command exists
	validCommands := []string{"start", "help", "status", "subscribe", "unsubscribe", "list"}
	if !cv.isValidCommand(ctx.Command, validCommands) {
		return errors.New("invalid command")
	}

	// Validate user permissions
	if !cv.hasPermission(ctx.UserID, ctx.Command) {
		return errors.New("insufficient permissions")
	}

	// Validate arguments based on command
	if err := cv.validateArguments(ctx.Command, ctx.Args); err != nil {
		return err
	}

	return nil
}

func (cv *CommandValidator) isValidCommand(command string, validCommands []string) bool {
	for _, valid := range validCommands {
		if command == valid {
			return true
		}
	}
	return false
}

func (cv *CommandValidator) hasPermission(userID int64, command string) bool {
	// For now, allow all users for basic commands
	basicCommands := []string{"start", "help", "status"}
	for _, basic := range basicCommands {
		if command == basic {
			return true
		}
	}

	// For advanced commands, check user permissions
	return cv.allowedUsers[userID]
}

func (cv *CommandValidator) validateArguments(command string, args []string) error {
	switch command {
	case "subscribe", "unsubscribe":
		if len(args) == 0 {
			return errors.New("project name is required")
		}
		if len(args[0]) < 2 {
			return errors.New("project name too short")
		}
	case "status":
		if len(args) > 1 {
			return errors.New("too many arguments for status command")
		}
	}
	return nil
}

// AddAllowedUser adds a user to the allowed users list
func (cv *CommandValidator) AddAllowedUser(userID int64) {
	cv.allowedUsers[userID] = true
}

// CommandRouter handles routing commands to appropriate handlers
type CommandRouter struct {
	handlers map[string]CommandHandler
}

// CommandHandler interface for command handlers
type CommandHandler interface {
	Handle(ctx *CommandContext) error
}

// NewCommandRouter creates a new command router
func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		handlers: make(map[string]CommandHandler),
	}
}

// RegisterHandler registers a command handler
func (cr *CommandRouter) RegisterHandler(command string, handler CommandHandler) {
	cr.handlers[strings.ToLower(command)] = handler
}

// RouteCommand routes a command to the appropriate handler
func (cr *CommandRouter) RouteCommand(ctx *CommandContext) error {
	handler, exists := cr.handlers[ctx.Command]
	if !exists {
		return errors.New("command handler not found")
	}

	return handler.Handle(ctx)
}

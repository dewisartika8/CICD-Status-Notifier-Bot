package telegram

import (
	"fmt"
	"strings"
)

// ResponseFormatter handles consistent response formatting for Telegram bot
type ResponseFormatter struct{}

func NewResponseFormatter() *ResponseFormatter {
	return &ResponseFormatter{}
}

// FormatError formats error messages with consistent styling
func (rf *ResponseFormatter) FormatError(title, message string) string {
	var response strings.Builder
	response.WriteString(fmt.Sprintf("❌ **%s**\n\n", title))
	response.WriteString(message)
	return response.String()
}

// FormatSuccess formats success messages with consistent styling
func (rf *ResponseFormatter) FormatSuccess(title, message string) string {
	var response strings.Builder
	response.WriteString(fmt.Sprintf("✅ **%s**\n\n", title))
	response.WriteString(message)
	return response.String()
}

// FormatInfo formats informational messages with consistent styling
func (rf *ResponseFormatter) FormatInfo(title, message string) string {
	var response strings.Builder
	response.WriteString(fmt.Sprintf("ℹ️ **%s**\n\n", title))
	response.WriteString(message)
	return response.String()
}

// FormatWarning formats warning messages with consistent styling
func (rf *ResponseFormatter) FormatWarning(title, message string) string {
	var response strings.Builder
	response.WriteString(fmt.Sprintf("⚠️ **%s**\n\n", title))
	response.WriteString(message)
	return response.String()
}

// FormatProjectNotFound formats project not found errors
func (rf *ResponseFormatter) FormatProjectNotFound(projectName string) string {
	return rf.FormatError("Project not found",
		fmt.Sprintf("Project `%s` is not being monitored.\n\n"+
			"Use `/projects` to see all monitored projects.", projectName))
}

// FormatServiceError formats service errors
func (rf *ResponseFormatter) FormatServiceError(operation string) string {
	return rf.FormatError(fmt.Sprintf("Error %s", operation),
		"Unable to complete the operation at the moment. Please try again later.")
}

// FormatInvalidCommand formats invalid command errors
func (rf *ResponseFormatter) FormatInvalidCommand(usage string) string {
	return rf.FormatError("Invalid command",
		fmt.Sprintf("Please check your command syntax.\n"+
			"Usage: %s", usage))
}

// Common error messages
const (
	ErrorProjectNotFound    = "Project not found or not monitored"
	ErrorServiceUnavailable = "Service temporarily unavailable"
	ErrorInvalidArguments   = "Invalid command arguments"
	ErrorPermissionDenied   = "Permission denied"
	ErrorInternalError      = "Internal error occurred"
)

// ErrorHandler handles different types of errors consistently
type ErrorHandler struct {
	formatter *ResponseFormatter
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		formatter: NewResponseFormatter(),
	}
}

// HandleProjectServiceError handles errors from project service calls
func (eh *ErrorHandler) HandleProjectServiceError(err error, operation string) string {
	if err == nil {
		return ""
	}

	// Check for specific error types (can be expanded based on actual domain errors)
	errorMsg := err.Error()

	switch {
	case strings.Contains(errorMsg, "not found"):
		return eh.formatter.FormatError("Project not found",
			"The requested project was not found or is not accessible.")
	case strings.Contains(errorMsg, "permission"):
		return eh.formatter.FormatError("Permission denied",
			"You don't have permission to access this resource.")
	case strings.Contains(errorMsg, "timeout"):
		return eh.formatter.FormatError("Request timeout",
			"The request took too long to complete. Please try again.")
	default:
		return eh.formatter.FormatServiceError(operation)
	}
}

// HandleValidationError handles validation errors for commands
func (eh *ErrorHandler) HandleValidationError(field, requirement string) string {
	return eh.formatter.FormatError("Validation error",
		fmt.Sprintf("Field `%s` %s", field, requirement))
}

package mcpsse

import "fmt"


type SSEClientCreationError struct {
	Message string
}

func (e *SSEClientCreationError) Error() string {
	return fmt.Sprintf("SSEClientCreationError: %s", e.Message)
}

type SSEClientStartError struct {
	Message string
}

func (e *SSEClientStartError) Error() string {
	return fmt.Sprintf("SSEClientStartError: %s", e.Message)
}

type SSEClientInitializationError struct {
	Message string
}

func (e *SSEClientInitializationError) Error() string {
	return fmt.Sprintf("SSEClientInitializationError: %s", e.Message)
}


type SSEGetToolsError struct {
	Message string
}

func (e *SSEGetToolsError) Error() string {
	return fmt.Sprintf("SSEGetToolsError: %s", e.Message)
}

type SSEToolCallError struct {
	Message string
}

func (e *SSEToolCallError) Error() string {
	return fmt.Sprintf("SSEToolCallError: %s", e.Message)
}

type SSEResultExtractionError struct {
	Message string
}

func (e *SSEResultExtractionError) Error() string {
	return fmt.Sprintf("SSEResultExtractionError: %s", e.Message)
}

type SSEListResourcesError struct {
	Message string
}

func (e *SSEListResourcesError) Error() string {
	return fmt.Sprintf("SSEListResourcesError: %s", e.Message)
}


type SSEReadResourceError struct {
	Message string
}

func (e *SSEReadResourceError) Error() string {
	return fmt.Sprintf("SSEReadResourceError: %s", e.Message)
}


type SSEListPromptsError struct {
	Message string
}

func (e *SSEListPromptsError) Error() string {
	return fmt.Sprintf("SSEListPromptsError: %s", e.Message)
}

type SSEGetPromptError struct {
	Message string
}

func (e *SSEGetPromptError) Error() string {
	return fmt.Sprintf("SSEGetPromptError: %s", e.Message)
}
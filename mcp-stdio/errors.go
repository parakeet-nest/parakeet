package mcpstdio

import "fmt"

type STDIOClientCreationError struct {
	Message string
}

func (e *STDIOClientCreationError) Error() string {
	return fmt.Sprintf("STDIOClientCreationError: %s", e.Message)
}

type STDIOClientInitializationError struct {
	Message string
}

func (e *STDIOClientInitializationError) Error() string {
	return fmt.Sprintf("STDIOClientInitializationError: %s", e.Message)
}

type STDIOGetToolsError struct {
	Message string
}

func (e *STDIOGetToolsError) Error() string {
	return fmt.Sprintf("STDIOGetToolsError: %s", e.Message)
}

type STDIOToolCallError struct {
	Message string
}

func (e *STDIOToolCallError) Error() string {
	return fmt.Sprintf("STDIOToolCallError: %s", e.Message)
}

type STDIOResultExtractionError struct {
	Message string
}

func (e *STDIOResultExtractionError) Error() string {
	return fmt.Sprintf("STDIOResultExtractionError: %s", e.Message)
}

type STDIOListResourcesError struct {
	Message string
}

func (e *STDIOListResourcesError) Error() string {
	return fmt.Sprintf("STDIOListResourcesError: %s", e.Message)
}

type STDIOReadResourceError struct {
	Message string
}

func (e *STDIOReadResourceError) Error() string {
	return fmt.Sprintf("STDIOReadResourceError: %s", e.Message)
}

type STDIOListPromptsError struct {
	Message string
}

func (e *STDIOListPromptsError) Error() string {
	return fmt.Sprintf("STDIOListPromptsError: %s", e.Message)
}

type STDIOGetPromptError struct {
	Message string
}

func (e *STDIOGetPromptError) Error() string {
	return fmt.Sprintf("STDIOGetPromptError: %s", e.Message)
}

package mcphttp

// HTTPClientCreationError represents an error that occurs during HTTP client creation
type HTTPClientCreationError struct {
	Message string
}

func (e *HTTPClientCreationError) Error() string {
	return e.Message
}

// HTTPClientStartError represents an error that occurs when starting the HTTP client
type HTTPClientStartError struct {
	Message string
}

func (e *HTTPClientStartError) Error() string {
	return e.Message
}

// HTTPClientInitializationError represents an error that occurs during HTTP client initialization
type HTTPClientInitializationError struct {
	Message string
}

func (e *HTTPClientInitializationError) Error() string {
	return e.Message
}

// HTTPGetToolsError represents an error that occurs when listing tools via HTTP
type HTTPGetToolsError struct {
	Message string
}

func (e *HTTPGetToolsError) Error() string {
	return e.Message
}

// HTTPToolCallError represents an error that occurs when calling a tool via HTTP
type HTTPToolCallError struct {
	Message string
}

func (e *HTTPToolCallError) Error() string {
	return e.Message
}

// HTTPResultExtractionError represents an error that occurs when extracting results from HTTP responses
type HTTPResultExtractionError struct {
	Message string
}

func (e *HTTPResultExtractionError) Error() string {
	return e.Message
}

// HTTPListResourcesError represents an error that occurs when listing resources via HTTP
type HTTPListResourcesError struct {
	Message string
}

func (e *HTTPListResourcesError) Error() string {
	return e.Message
}

// HTTPReadResourceError represents an error that occurs when reading a resource via HTTP
type HTTPReadResourceError struct {
	Message string
}

func (e *HTTPReadResourceError) Error() string {
	return e.Message
}

// HTTPListPromptsError represents an error that occurs when listing prompts via HTTP
type HTTPListPromptsError struct {
	Message string
}

func (e *HTTPListPromptsError) Error() string {
	return e.Message
}

// HTTPGetPromptError represents an error that occurs when getting a prompt via HTTP
type HTTPGetPromptError struct {
	Message string
}

func (e *HTTPGetPromptError) Error() string {
	return e.Message
}
package completion

import "fmt"

type ModelNotFoundError struct {
	Code    int
	Message string
	Model   string
}

func (e *ModelNotFoundError) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Model: %s", e.Code, e.Message, e.Model)
}

type NoSuchOllamaHostError struct {
	Host string
	Message string
}

func (e *NoSuchOllamaHostError) Error() string {
	return fmt.Sprintf("Host: %s, Message: %s", e.Host, e.Message)
}



type CompletionError struct {
	Error string
}

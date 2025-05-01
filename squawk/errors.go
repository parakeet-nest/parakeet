package squawk

// LastError gets or sets the most recent error encountered during processing.
// When called without parameters, it returns the last error. When called with
// an error parameter, it sets the last error and returns it.
//
// Parameters:
//   - optionalError: Optional error to set as the last error
//
// Returns:
//   - error: The current last error or the newly set error
//
// Example with error handling:
//   squawk := New().
//     Model("mistral:latest").
//     Provider(provider.Ollama).
//     System("You are a Go expert").
//     Chat(func(answer llm.Answer, self *Squawk, err error) {
//         if err != nil {
//             self.LastError(err)  // Set the error
//             return
//         }
//     })
//   
//   if err := squawk.LastError(); err != nil {
//     fmt.Printf("Last operation failed: %v\n", err)
//   }
//
// Example with error tracking:
//   squawk := New().
//     Model("codellama:13b").
//     Provider(provider.Ollama).
//     ChatStream(func(answer llm.Answer, self *Squawk) error {
//         if answer.Error != nil {
//             self.LastError(answer.Error)  // Track streaming error
//             return answer.Error
//         }
//         return nil
//     })
func (s *Squawk) LastError(optionalError ...error) error {
	if len(optionalError) > 0 {
		s.lastError = optionalError[0]
	}
	return s.lastError
}

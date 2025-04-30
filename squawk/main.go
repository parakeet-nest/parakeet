package squawk

import (
	"log"

	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"
)

type Squawk struct {
	setOfMessages []llm.Message
	baseUrl       string
	apiUrl        string
	provider      string
	model         string
	options       llm.Options
	openAPIKey    string
	lastAnswer    llm.Answer
	lastError	 error
}

/*
I want to create a DSL for parakeet to make it easier to use.


*/

func New(model string) *Squawk {
	s := &Squawk{
		setOfMessages: []llm.Message{},
		model:         model,
		baseUrl:       "",
		apiUrl:        "",
		provider:      provider.Ollama,
		options:       llm.Options{},
		lastAnswer:   llm.Answer{},
		lastError:    nil,
	}
	// Initialize the Squawk instance with the provided model and arguments
	// You can add any necessary initialization logic here

	// The possible value for args[0] is a string, which is the llm provider

	return s
}

func (s *Squawk) System(message string, optionalParameters ...string) *Squawk {
	if len(optionalParameters) > 0 {
		s.setOfMessages = append(
			s.setOfMessages,
			llm.Message{Role: "system", Content: message, Label: optionalParameters[0]},
		)

	} else {
		s.setOfMessages = append(s.setOfMessages, llm.Message{Role: "system", Content: message})
	}

	return s
}

func (s *Squawk) User(message string, optionalParameters ...string) *Squawk {
	if len(optionalParameters) > 0 {
		s.setOfMessages = append(
			s.setOfMessages,
			llm.Message{Role: "user", Content: message, Label: optionalParameters[0]},
		)
	} else {
		s.setOfMessages = append(s.setOfMessages, llm.Message{Role: "user", Content: message})
	}
	return s
}

func (s *Squawk) Assistant(message string, optionalParameters ...string) *Squawk {
	if len(optionalParameters) > 0 {
		s.setOfMessages = append(
			s.setOfMessages,
			llm.Message{Role: "assistant", Content: message, Label: optionalParameters[0]},
		)
	} else {
		s.setOfMessages = append(s.setOfMessages, llm.Message{Role: "assistant", Content: message})
	}
	return s
}

func (s *Squawk) SetOfMessages(messages ...interface{}) *Squawk {
	// SetOfMessages creates or extends a conversation with provided messages
	// It can accept either single messages or slices of messages as variadic parameters
	for _, msg := range messages {
		switch m := msg.(type) {
		case llm.Message:
			s.setOfMessages = append(s.setOfMessages, m)
		case []llm.Message:
			s.setOfMessages = append(s.setOfMessages, m...)
		}
	}
	return s
}

func (s *Squawk) BaseURL(url string) *Squawk {
	s.baseUrl = url
	//log.Println("🟢", url, s.baseUrl)
	return s
}

func (s *Squawk) Provider(llmProvider string, parameters ...string) *Squawk {
	s.provider = llmProvider
	switch llmProvider {
	case provider.Ollama:
		//log.Println("🦙", llmProvider)
		if s.baseUrl == "" {
			s.apiUrl = "http://localhost:11434"
		} else {
			s.apiUrl = s.baseUrl
		}
	case provider.DockerModelRunner:
		//log.Println("🐳", llmProvider)
		if s.baseUrl == "" {
			s.apiUrl = "http://localhost:12434/engines/llama.cpp/v1"
		} else {
			s.apiUrl = s.baseUrl + "/engines/llama.cpp/v1"
		}

	case provider.OpenAI:
		log.Println("🔵", llmProvider, parameters[0])
		s.apiUrl = "https://api.openai.com/v1"
		s.openAPIKey = parameters[0]
	default: // Ollama
		s.apiUrl = s.baseUrl
	}
	//log.Println("🔴", s.apiUrl)
	return s
}

// You can change of model
func (s *Squawk) Model(model string) *Squawk {
	s.model = model
	return s
}

// TODO: simplify this part
func (s *Squawk) Options(options llm.Options) *Squawk {
	s.options = options
	return s
}

func (s *Squawk) chatExec() (llm.Answer, error) {
	query := llm.Query{
		Model:    s.model,
		Messages: s.setOfMessages,
		Options:  s.options,
	}
	// Call the chat function with the provided model and options
	//log.Println("🔵", s.apiUrl)

	answer, err := completion.Chat(s.apiUrl, query, s.provider, s.openAPIKey)
	if err != nil {
		return llm.Answer{}, err
	}
	s.lastAnswer = answer
	return answer, nil
}

func (s *Squawk) chatStreamExec(callBack func(answer llm.Answer) error) (llm.Answer, error) {
	query := llm.Query{
		Model:    s.model,
		Messages: s.setOfMessages,
		Options:  s.options,
	}
	answer, err := completion.ChatStream(s.apiUrl, query, callBack , s.provider, s.openAPIKey)
	if err != nil {
		return llm.Answer{}, err
	}
	s.lastAnswer = answer
	return answer, nil

}

func (s *Squawk) Chat(callBack func(answer llm.Answer, self *Squawk, err error)) *Squawk {
	answer, err := s.chatExec()
	if err != nil {
		callBack(answer, s, err)
		s.lastError = err
		return s
	}
	callBack(answer, s, nil)
	//s.lastAnswer = answer
	return s
}

func (s *Squawk) ChatStream(callBack func(answer llm.Answer, self *Squawk) error) *Squawk {
	answer, err := s.chatStreamExec(func(answer llm.Answer) error {
		return callBack(answer, s)
	})
	if err != nil {
		s.lastAnswer = llm.Answer{}
		s.lastError = err
		return s
	} 
	s.lastAnswer = answer
	s.lastError = nil
	return s
}

func (s *Squawk) Exec(callBack func(self *Squawk)) *Squawk {
	callBack(s)
	return s
}

func (s *Squawk) SaveAnswer(optionalParameters ...string) *Squawk {
	if len(optionalParameters) > 0 {
		s.setOfMessages = append(
			s.setOfMessages,
			llm.Message{Role: "assistant", Content: s.lastAnswer.Message.Content, Label: optionalParameters[0]},
		)
	} else {
		s.setOfMessages = append(s.setOfMessages, llm.Message{Role: "assistant", Content: s.lastAnswer.Message.Content})
	}
	return s

}

func (s *Squawk) LastAnswer(optionalAnswer ...llm.Answer) llm.Answer {
	if len(optionalAnswer) > 0 {
		s.lastAnswer = optionalAnswer[0]
	} 
	return s.lastAnswer
}
func (s *Squawk) LastError(optionalError ...error) error {
	if len(optionalError) > 0 {
		s.lastError = optionalError[0]
	}
	return s.lastError
}

// TODO: remove messages by label
// TODO: replace messages by label
// TODO: json format output -> Structured output with OpenAI???
// TODO: tools (+ MCP conversion)
	// https://platform.openai.com/docs/guides/structured-outputs?api-mode=chat&lang=curl
	// response_format
// TODO: embeddings
// TODO: similarity search

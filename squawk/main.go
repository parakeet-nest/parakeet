package squawk

/*
"Squawk is the jQuery of generative AI"
Squawk plays a similar role in generative AI that jQuery did for JavaScript.
It simplifies common tasks, making the technology more accessible and easier to work with.
*/

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/embeddings"
	"github.com/parakeet-nest/parakeet/enums/provider"
	"github.com/parakeet-nest/parakeet/llm"
)

type Squawk struct {
	setOfMessages   []llm.Message
	baseUrl         string
	apiUrl          string
	provider        string
	chatModel       string
	embeddingsModel string
	options         llm.Options
	openAPIKey      string
	lastAnswer      llm.Answer
	lastError       error
	schema          map[string]any // for structured output

	// embeddings
	vectorStore  embeddings.VectorStore
	similarities []llm.VectorRecord
}

func New() *Squawk {
	s := &Squawk{
		setOfMessages:   []llm.Message{},
		chatModel:       "",
		embeddingsModel: "",
		baseUrl:         "",
		apiUrl:          "",
		provider:        provider.Ollama,
		options:         llm.Options{},
		lastAnswer:      llm.Answer{},
		lastError:       nil,

		vectorStore:  nil,
		similarities: []llm.VectorRecord{},
	}
	// Initialize the Squawk instance with the provided model and arguments
	// You can add any necessary initialization logic here

	// The possible value for args[0] is a string, which is the llm provider

	return s
}

// You can change of model
func (s *Squawk) Model(model string) *Squawk {
	s.chatModel = model
	return s
}
func (s *Squawk) EmbeddingsModel(model string) *Squawk {
	s.embeddingsModel = model
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

// Provider sets the LLM (Large Language Model) provider for the Squawk instance
// and configures the API URL and other parameters based on the selected provider.
//
// Parameters:
//   - llmProvider: A string representing the LLM provider. Supported values include:
//   - provider.Ollama: Configures the API URL for the Ollama provider.
//   - provider.DockerModelRunner: Configures the API URL for the Docker Model Runner provider.
//   - provider.OpenAI: Configures the API URL for OpenAI and sets the OpenAI API key.
//   - parameters: Optional additional parameters. For provider.OpenAI, the first parameter
//     should be the OpenAI API key.
//
// Behavior:
//   - If the baseUrl is not set, a default API URL is assigned based on the provider.
//   - For provider.OpenAI, the first parameter in the `parameters` slice is used as the API key.
//   - If an unsupported provider is specified, the baseUrl is used as the API URL.
//
// Returns:
//   - A pointer to the updated Squawk instance.
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
		//log.Println("🔵", llmProvider, parameters[0])
		if s.baseUrl == "" {
			s.apiUrl = "https://api.openai.com/v1"
		} else {
			s.apiUrl = s.baseUrl + "/v1"
		}
		s.openAPIKey = parameters[0]
	default: // Ollama
		s.apiUrl = s.baseUrl
	}
	return s
}

// TODO: simplify this part
func (s *Squawk) Options(options llm.Options) *Squawk {
	s.options = options
	return s
}

func (s *Squawk) Schema(schema map[string]any) *Squawk {
	s.schema = schema
	return s
}

func (s *Squawk) chatExec() (llm.Answer, error) {
	query := llm.Query{
		Model:    s.chatModel,
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

func (s *Squawk) structuredOutputExec() (llm.Answer, error) {
	query := llm.Query{
		Model:    s.chatModel,
		Messages: s.setOfMessages,
		Options:  s.options,
		Format:   s.schema,
		Raw:      false,
	}
	answer, err := completion.Chat(s.apiUrl, query, s.provider, s.openAPIKey)
	if err != nil {
		return llm.Answer{}, err
	}
	s.lastAnswer = answer
	return answer, nil
}

func (s *Squawk) StructuredOutput(callBack func(answer llm.Answer, self *Squawk, err error)) *Squawk {
	answer, err := s.structuredOutputExec()
	if err != nil {
		callBack(answer, s, err)
		s.lastError = err
		return s
	}
	callBack(answer, s, nil)
	//s.lastAnswer = answer
	return s
}

func (s *Squawk) chatStreamExec(callBack func(answer llm.Answer) error) (llm.Answer, error) {
	query := llm.Query{
		Model:    s.chatModel,
		Messages: s.setOfMessages,
		Options:  s.options,
	}
	answer, err := completion.ChatStream(s.apiUrl, query, callBack, s.provider, s.openAPIKey)
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

func (s *Squawk) Cmd(callBack func(self *Squawk)) *Squawk {
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

func (s *Squawk) RemoveMessageByLabel(label string) *Squawk {
	// Remove messages by label
	var newMessages []llm.Message
	for _, message := range s.setOfMessages {
		if message.Label != label {
			newMessages = append(newMessages, message)
		}
	}
	s.setOfMessages = newMessages
	return s
}

// TODO: tools (+ MCP conversion)
// https://platform.openai.com/docs/guides/structured-outputs?api-mode=chat&lang=curl
// response_format
// ...
//TODO: add default options ()

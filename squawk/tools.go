package squawk

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/llm"
)

func (s *Squawk) Tools(toolsList []llm.Tool) *Squawk {

	s.tools = toolsList
	return s
}

func (s *Squawk) functionCallingExec() (llm.Answer, error) {
	query := llm.Query{
		Model:    s.chatModel,
		Messages: s.setOfMessages,
		Tools:    s.tools,
		Options:  s.options,
	}

	answer, err := completion.Chat(s.apiUrl, query, s.provider, s.openAPIKey)
	if err != nil {
		return llm.Answer{}, err
	}
	s.lastAnswer = answer
	s.toolCalls = answer.Message.ToolCalls.Tools()
	return answer, nil
}

func (s *Squawk) FunctionCalling(callBack func(answer llm.Answer, self *Squawk, err error)) *Squawk {

	answer, err := s.functionCallingExec()
	if err != nil {
		callBack(answer, s, err)
		s.lastError = err
		return s
	}
	callBack(answer, s, nil)
	//s.lastAnswer = answer
	return s
}

func (s *Squawk) ToolCalls() []llm.ToolCall {
	return s.toolCalls
}

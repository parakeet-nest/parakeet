package squawk

import (
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/prompt"
)

func (s *Squawk) ForKids(message string, optionalParameters ...string) *Squawk {
	message = prompt.ForKids(message)
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

func (s *Squawk) Brief(message string, optionalParameters ...string) *Squawk {
	message = prompt.Brief(message)
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


func (s *Squawk) AsAPoem(message string, optionalParameters ...string) *Squawk {
	message = prompt.AsAPoem(message)
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

func (s *Squawk) AdvantagesOnly(message string, optionalParameters ...string) *Squawk {
	message = prompt.AdvantagesOnly(message)
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

func (s *Squawk) AsARecipe(message string, optionalParameters ...string) *Squawk {
	message = prompt.AsARecipe(message)
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

func (s *Squawk) Timeline(message string, optionalParameters ...string) *Squawk {
	message = prompt.Timeline(message)
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

func (s *Squawk) Comparison(message string, optionalParameters ...string) *Squawk {
	message = prompt.Comparison(message)
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

func (s *Squawk) Opinion(message string, optionalParameters ...string) *Squawk {
	message = prompt.Opinion(message)
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

func (s *Squawk) Factual(message string, optionalParameters ...string) *Squawk {
	message = prompt.Factual(message)
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

func (s *Squawk) StepByStep(message string, optionalParameters ...string) *Squawk {
	message = prompt.StepByStep(message)
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

func (s *Squawk) ProsAndCons(message string, optionalParameters ...string) *Squawk {
	message = prompt.ProsAndCons(message)
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

func (s *Squawk) AsAStory(message string, optionalParameters ...string) *Squawk {
	message = prompt.AsAStory(message)
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

func (s *Squawk) InLaymansTerms(message string, optionalParameters ...string) *Squawk {
	message = prompt.InLaymansTerms(message)
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

func (s *Squawk) Summarize(message string, optionalParameters ...string) *Squawk {
	message = "Summarize the following text: " + message
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

// ???: is this useful?
func (s *Squawk) SummarizeLastAnswer(optionalParameters ...string) *Squawk {
	message := "Summarize the following text: " + s.lastAnswer.Message.Content
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
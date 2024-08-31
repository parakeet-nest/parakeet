package data

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

type QA struct {
	Question string   `json:"question" yaml:"question"`
	Answer   string   `json:"answer" yaml:"answer"`
	MetaData string   `json:"metadata" yaml:"metadata"`
	KeyWords []string `json:"keywords" yaml:"keywords"`
}

func (qa *QA) ToText() string {
	if len(qa.KeyWords) == 0 {
		return fmt.Sprintf("QUESTION: %s\nANSWER: %s\nMETADATA: %s", qa.Question, qa.Answer, qa.MetaData)
	}
	if qa.MetaData == "" {
		return fmt.Sprintf("QUESTION: %s\nANSWER: %s\nKEYWORDS: %s", qa.Question, qa.Answer, strings.Join(qa.KeyWords, ", "))
	}
	if len(qa.KeyWords) == 0 && qa.MetaData == "" {
		return fmt.Sprintf("QUESTION: %s\nANSWER: %s", qa.Question, qa.Answer)
	}
	return fmt.Sprintf("QUESTION: %s\nANSWER: %s\nMETADATA: %s\nKEYWORDS: %s", qa.Question, qa.Answer, qa.MetaData, strings.Join(qa.KeyWords, ", "))
}

func (qa *QA) ToJson() string {
	jsonBytes, err := json.MarshalIndent(qa, "", "  ")
	if err != nil {
		// handle error
		return "{}"
	}
	return string(jsonBytes)
}

func (qa *QA) FromJson(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), qa)
}

func QAListFromJson(jsonString string) ([]QA, error) {
	var qaList []QA
	err := json.Unmarshal([]byte(jsonString), &qaList)
	return qaList, err
}

func QAListToJson(qaList []QA) string {
	jsonBytes, err := json.MarshalIndent(qaList, "", "  ")
	if err != nil {
		// handle error
		return "[]"
	}
	return string(jsonBytes)
}

func QAListToText(qaList []QA) string {
	var qaTexts []string
	for _, qa := range qaList {
		qaTexts = append(qaTexts, qa.ToText())
	}
	return strings.Join(qaTexts, "\n\n")
}

func QAListFromYaml(yamlString string) ([]QA, error) {
	var qaList []QA
	err := yaml.Unmarshal([]byte(yamlString), &qaList)
	return qaList, err
}

func QAListToYaml(qaList []QA) string {
	yamlBytes, err := yaml.Marshal(qaList)
	if err != nil {
		// handle error
		return "[]"
	}
	return string(yamlBytes)
}

type IO struct {
	Input    string   `json:"input" yaml:"input"`
	Output   string   `json:"output" yaml:"output"`
	MetaData string   `json:"metadata" yaml:"metadata"`
	KeyWords []string `json:"keywords" yaml:"keywords"`
}

func (io *IO) ToText() string {

	if len(io.KeyWords) == 0 {
		return fmt.Sprintf("INPUT: %s\nOUTPUT: %s\nMETADATA: %s", io.Input, io.Output, io.MetaData)
	}
	if io.MetaData == "" {
		return fmt.Sprintf("INPUT: %s\nOUTPUT: %s\nKEYWORDS: %s", io.Input, io.Output, strings.Join(io.KeyWords, ", "))
	}
	if len(io.KeyWords) == 0 && io.MetaData == "" {
		return fmt.Sprintf("INPUT: %s\nOUTPUT: %s", io.Input, io.Output)
	}
	return fmt.Sprintf("INPUT: %s\nOUTPUT: %s\nMETADATA: %s\nKEYWORDS: %s", io.Input, io.Output, io.MetaData, strings.Join(io.KeyWords, ", "))
}

func (io *IO) ToJson() string {
	jsonBytes, err := json.MarshalIndent(io, "", "  ")
	if err != nil {
		// handle error
		return "{}"
	}
	return string(jsonBytes)
}

func (io *IO) FromJson(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), io)
}

func IOListFromJson(jsonString string) ([]IO, error) {
	var ioList []IO
	err := json.Unmarshal([]byte(jsonString), &ioList)
	return ioList, err
}

func IOListToJson(ioList []IO) string {
	jsonBytes, err := json.MarshalIndent(ioList, "", "  ")
	if err != nil {
		// handle error
		return "[]"
	}
	return string(jsonBytes)
}

func IOListToText(ioList []IO) string {
	var ioTexts []string
	for _, io := range ioList {
		ioTexts = append(ioTexts, io.ToText())
	}
	return strings.Join(ioTexts, "\n\n")
}

func IOListFromYaml(yamlString string) ([]IO, error) {
	var ioList []IO
	err := yaml.Unmarshal([]byte(yamlString), &ioList)
	return ioList, err
}

func IOListToYaml(ioList []IO) string {
	yamlBytes, err := yaml.Marshal(ioList)
	if err != nil {
		// handle error
		return "[]"
	}
	return string(yamlBytes)
}


type Card struct {
	Title    string   `json:"title" yaml:"title"`
	Body     string   `json:"body" yaml:"body"`
	MetaData string   `json:"metadata" yaml:"metadata"`
	KeyWords []string `json:"keywords" yaml:"keywords"`
}

func (card *Card) ToText() string {
	if len(card.KeyWords) == 0 {
		return fmt.Sprintf("TITLE: %s\nBODY: %s\nMETADATA: %s", card.Title, card.Body, card.MetaData)
	}
	if card.MetaData == "" {
		return fmt.Sprintf("TITLE: %s\nBODY: %s\nKEYWORDS: %s", card.Title, card.Body, strings.Join(card.KeyWords, ", "))
	}
	if len(card.KeyWords) == 0 && card.MetaData == "" {
		return fmt.Sprintf("TITLE: %s\nBODY: %s", card.Title, card.Body)
	}
	return fmt.Sprintf("TITLE: %s\nBODY: %s\nMETADATA: %s\nKEYWORDS: %s", card.Title, card.Body, card.MetaData, strings.Join(card.KeyWords, ", "))
}

func (card *Card) ToJson() string {
	jsonBytes, err := json.MarshalIndent(card, "", "  ")
	if err != nil {
		// handle error
		return "{}"
	}
	return string(jsonBytes)
}

func (card *Card) FromJson(jsonString string) error {
	return json.Unmarshal([]byte(jsonString), card)
}

func CardListFromJson(jsonString string) ([]Card, error) {
	var cardList []Card
	err := json.Unmarshal([]byte(jsonString), &cardList)
	return cardList, err
}

func CardListToJson(cardList []Card) string {
	jsonBytes, err := json.MarshalIndent(cardList, "", "  ")
	if err != nil {
		// handle error
		return "[]"
	}
	return string(jsonBytes)
}

func CardListToText(cardList []Card) string {
	var cardTexts []string
	for _, card := range cardList {
		cardTexts = append(cardTexts, card.ToText())
	}
	return strings.Join(cardTexts, "\n\n")
}

func CardListFromYaml(yamlString string) ([]Card, error) {
	var cardList []Card
	err := yaml.Unmarshal([]byte(yamlString), &cardList)
	return cardList, err
}

func CardListToYaml(cardList []Card) string {
	yamlBytes, err := yaml.Marshal(cardList)
	if err != nil {
		// handle error
		return "[]"
	}
	return string(yamlBytes)
}
package testing_data

import (
	"fmt"
	"testing"
	"github.com/parakeet-nest/parakeet/data"
)

func TestQAToText(t *testing.T) {
	qa := data.QA{
		Question: "What is the meaning of life?",
		Answer:   "42",
		MetaData: "Don't panic",
		KeyWords: []string{"life", "meaning", "universe"},
	}

	res := qa.ToText()

	expected := `QUESTION: What is the meaning of life?
ANSWER: 42
METADATA: Don't panic
KEYWORDS: life, meaning, universe`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)
}


func TestQAToJson(t *testing.T) {
	qa := data.QA{
		Question: "What is the meaning of life?",
		Answer:   "42",
		MetaData: "Don't panic",
		KeyWords: []string{"life", "meaning", "universe"},
	}

	res := qa.ToJson()

	expected := `{
  "question": "What is the meaning of life?",
  "answer": "42",
  "metadata": "Don't panic",
  "keywords": [
    "life",
    "meaning",
    "universe"
  ]
}`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)
	
}

func TestQAFromJson(t *testing.T) {
	jsonString := `{
  "question": "What is the meaning of life?",
  "answer": "42",
  "metadata": "Don't panic",
  "keywords": [
	"life",
	"meaning",
	"universe"
  ]
}`

	qa := data.QA{}
	err := qa.FromJson(jsonString)

	if err != nil {
		t.Fatal(err)
	}

	if qa.Question != "What is the meaning of life?" {
		t.Fatalf("unexpected question: %s", qa.Question)
	}

	fmt.Println("✅", qa.Question)
}

func TestQAListFromJson(t *testing.T) {
	jsonString := `[
  {
	"question": "What is the meaning of life?",
	"answer": "42",
	"metadata": "Don't panic",
	"keywords": [
	  "life",
	  "meaning",
	  "universe"
	]
  },
  {
	"question": "What is the answer to the ultimate question of life, the universe, and everything?",
	"answer": "42",
	"metadata": "Don't panic",
	"keywords": [
	  "life",
	  "meaning",
	  "universe"
	]
  }
]`

	qaList, err := data.QAListFromJson(jsonString)

	if err != nil {
		t.Fatal(err)
	}

	if len(qaList) != 2 {
		t.Fatalf("unexpected length: %d", len(qaList))
	}

	fmt.Println("✅", len(qaList))
}

func TestQAListToJson(t *testing.T) {
	qaList := []data.QA{
		{
			Question: "What is the meaning of life?",
			Answer:   "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
		{
			Question: "What is the answer to the ultimate question of life, the universe, and everything?",
			Answer:   "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
	}

	res := data.QAListToJson(qaList)

	expected := `[
  {
    "question": "What is the meaning of life?",
    "answer": "42",
    "metadata": "Don't panic",
    "keywords": [
      "life",
      "meaning",
      "universe"
    ]
  },
  {
    "question": "What is the answer to the ultimate question of life, the universe, and everything?",
    "answer": "42",
    "metadata": "Don't panic",
    "keywords": [
      "life",
      "meaning",
      "universe"
    ]
  }
]`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)
}

func TestQAListToText(t *testing.T) {
	qaList := []data.QA{
		{
			Question: "What is the meaning of life?",
			Answer:   "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
		{
			Question: "What is the answer to the ultimate question of life, the universe, and everything?",
			Answer:   "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
	}

	res := data.QAListToText(qaList)

	expected := `QUESTION: What is the meaning of life?
ANSWER: 42
METADATA: Don't panic
KEYWORDS: life, meaning, universe

QUESTION: What is the answer to the ultimate question of life, the universe, and everything?
ANSWER: 42
METADATA: Don't panic
KEYWORDS: life, meaning, universe`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)

}


func TestQAFromYaml(t *testing.T) {
	yamlString := `- question: What is the meaning of life?
  answer: "42"
  metadata: Don't panic
  keywords:
  - life
  - meaning
  - universe
- question: What is the answer to the ultimate question of life, the universe, and
    everything?
  answer: "42"
  metadata: Don't panic
  keywords:
  - life
  - meaning
  - universe
`

	qaList, err := data.QAListFromYaml(yamlString)

	if err != nil {
		t.Fatal(err)
	}	

	if len(qaList) != 2 {
		t.Fatalf("unexpected length: %d", len(qaList))
	}

	fmt.Println("✅", len(qaList))
}


func TestQAListToYaml(t *testing.T) {
	qaList := []data.QA{
		{
			Question: "What is the meaning of life?",
			Answer:   "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
		{
			Question: "What is the answer to the ultimate question of life, the universe, and everything?",
			Answer:   "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
	}

	res := data.QAListToYaml(qaList)


	expected := `- question: What is the meaning of life?
  answer: "42"
  metadata: Don't panic
  keywords:
  - life
  - meaning
  - universe
- question: What is the answer to the ultimate question of life, the universe, and
    everything?
  answer: "42"
  metadata: Don't panic
  keywords:
  - life
  - meaning
  - universe
`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)
}


func TestIOToText(t *testing.T) {
	io := data.IO{
		Input:    "What is the meaning of life?",
		Output:   "42",
		MetaData: "Don't panic",
		KeyWords: []string{"life", "meaning", "universe"},
	}
	
	res := io.ToText()

	expected := `INPUT: What is the meaning of life?
OUTPUT: 42
METADATA: Don't panic
KEYWORDS: life, meaning, universe`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)

}

func TestIOToJson(t *testing.T) {
	io := data.IO{
		Input:    "What is the meaning of life?",
		Output:   "42",
		MetaData: "Don't panic",
		KeyWords: []string{"life", "meaning", "universe"},
	}
	
	res := io.ToJson()

	expected := `{
  "input": "What is the meaning of life?",
  "output": "42",
  "metadata": "Don't panic",
  "keywords": [
    "life",
    "meaning",
    "universe"
  ]
}`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)

}

func TestIOFromJson(t *testing.T) {
	jsonString := `{
  "input": "What is the meaning of life?",
  "output": "42",
  "metadata": "Don't panic",
  "keywords": [
	"life",
	"meaning",
	"universe"
  ]
}`

	io := data.IO{}
	err := io.FromJson(jsonString)

	if err != nil {
		t.Fatal(err)
	}

	if io.Input != "What is the meaning of life?" {
		t.Fatalf("unexpected input: %s", io.Input)
	}

	fmt.Println("✅", io.Input)
}

func TestIOListFromJson(t *testing.T) {	
	jsonString := `[
  {
	"input": "What is the meaning of life?",
	"output": "42",
	"metadata": "Don't panic",
	"keywords": [
	  "life",
	  "meaning",
	  "universe"
	]
  },
  {
	"input": "What is the answer to the ultimate question of life, the universe, and everything?",
	"output": "42",
	"metadata": "Don't panic",
	"keywords": [
	  "life",
	  "meaning",
	  "universe"
	]
  }
]`

	ioList, err := data.IOListFromJson(jsonString)

	if err != nil {
		t.Fatal(err)
	}

	if len(ioList) != 2 {
		t.Fatalf("unexpected length: %d", len(ioList))
	}

	fmt.Println("✅", len(ioList))
}


func TestIOListToJson(t *testing.T) {
	ioList := []data.IO{
		{
			Input:    "What is the meaning of life?",
			Output:   "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
		{
			Input:    "What is the answer to the ultimate question of life, the universe, and everything?",
			Output:   "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
	}

	res := data.IOListToJson(ioList)

	expected := `[
  {
    "input": "What is the meaning of life?",
    "output": "42",
    "metadata": "Don't panic",
    "keywords": [
      "life",
      "meaning",
      "universe"
    ]
  },
  {
    "input": "What is the answer to the ultimate question of life, the universe, and everything?",
    "output": "42",
    "metadata": "Don't panic",
    "keywords": [
      "life",
      "meaning",
      "universe"
    ]
  }
]`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)
}

func TestIOListToText(t *testing.T) {
	ioList := []data.IO{
		{
			Input:    "What is the meaning of life?",
			Output:   "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
		{
			Input:    "What is the answer to the ultimate question of life, the universe, and everything?",
			Output:   "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
	}

	res := data.IOListToText(ioList)

	expected := `INPUT: What is the meaning of life?
OUTPUT: 42
METADATA: Don't panic
KEYWORDS: life, meaning, universe

INPUT: What is the answer to the ultimate question of life, the universe, and everything?
OUTPUT: 42
METADATA: Don't panic
KEYWORDS: life, meaning, universe`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)

}

func TestIOFromYaml(t *testing.T) {
	yamlString := `- input: What is the meaning of life?
  output: "42"
  metadata: Don't panic
  keywords:
  - life
  - meaning
  - universe
- input: What is the answer to the ultimate question of life, the universe, and everything?
  output: "42"
  metadata: Don't panic
  keywords:
  - life
  - meaning
  - universe
`

	ioList, err := data.IOListFromYaml(yamlString)
	
	if err != nil {
		t.Fatal(err)
	}

	if len(ioList) != 2 {
		t.Fatalf("unexpected length: %d", len(ioList))
	}
	fmt.Println("✅", len(ioList))
}

	
func TestIOListToYaml(t *testing.T) {
	ioList := []data.IO{
		{
			Input:    "What is the meaning of life?",
			Output:   "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
		{
			Input:    "What is the answer to the ultimate question of life, the universe, and everything?",
			Output:   "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
	}

	res := data.IOListToYaml(ioList)

	expected := `- input: What is the meaning of life?
  output: "42"
  metadata: Don't panic
  keywords:
  - life
  - meaning
  - universe
- input: What is the answer to the ultimate question of life, the universe, and everything?
  output: "42"
  metadata: Don't panic
  keywords:
  - life
  - meaning
  - universe
`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)
}

func TestCardToText(t *testing.T) {
	card := data.Card{
		Title: "What is the meaning of life?",
		Body:  "42",
		MetaData: "Don't panic",
		KeyWords: []string{"life", "meaning", "universe"},
	}
	
	res := card.ToText()	

	expected := `TITLE: What is the meaning of life?
BODY: 42
METADATA: Don't panic
KEYWORDS: life, meaning, universe`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)
}

func TestCardToJson(t *testing.T) {
	card := data.Card{
		Title: "What is the meaning of life?",
		Body:  "42",
		MetaData: "Don't panic",
		KeyWords: []string{"life", "meaning", "universe"},
	}
	
	res := card.ToJson()

	expected := `{
  "title": "What is the meaning of life?",
  "body": "42",
  "metadata": "Don't panic",
  "keywords": [
    "life",
    "meaning",
    "universe"
  ]
}`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)
}

func TestCardFromJson(t *testing.T) {
	jsonString := `{
  "title": "What is the meaning of life?",
  "body": "42",
  "metadata": "Don't panic",
  "keywords": [
	"life",
	"meaning",
	"universe"
  ]
}`

	card := data.Card{}
	err := card.FromJson(jsonString)

	if err != nil {
		t.Fatal(err)
	}

	if card.Title != "What is the meaning of life?" {
		t.Fatalf("unexpected title: %s", card.Title)
	}

	fmt.Println("✅", card.Title)
}

func TestCardListFromJson(t *testing.T) {
	jsonString := `[
  {
	"title": "What is the meaning of life?",
	"body": "42",
	"metadata": "Don't panic",
	"keywords": [
	  "life",
	  "meaning",
	  "universe"
	]
  },
  {
	"title": "What is the answer to the ultimate question of life, the universe, and everything?",
	"body": "42",
	"metadata": "Don't panic",
	"keywords": [
	  "life",
	  "meaning",
	  "universe"
	]
  }
]`

	cardList, err := data.CardListFromJson(jsonString)

	if err != nil {
		t.Fatal(err)
	}

	if len(cardList) != 2 {
		t.Fatalf("unexpected length: %d", len(cardList))
	}

	fmt.Println("✅", len(cardList))
}

func TestCardListToJson(t *testing.T) {
	cardList := []data.Card{
		{
			Title: "What is the meaning of life?",
			Body:  "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
		{
			Title: "What is the answer to the ultimate question of life, the universe, and everything?",
			Body:  "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
	}

	res := data.CardListToJson(cardList)

	expected := `[
  {
    "title": "What is the meaning of life?",
    "body": "42",
    "metadata": "Don't panic",
    "keywords": [
      "life",
      "meaning",
      "universe"
    ]
  },
  {
    "title": "What is the answer to the ultimate question of life, the universe, and everything?",
    "body": "42",
    "metadata": "Don't panic",
    "keywords": [
      "life",
      "meaning",
      "universe"
    ]
  }
]`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)
}

func TestCardListToText(t *testing.T) {
	cardList := []data.Card{
		{
			Title: "What is the meaning of life?",
			Body:  "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
		{
			Title: "What is the answer to the ultimate question of life, the universe, and everything?",
			Body:  "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
	}

	res := data.CardListToText(cardList)

	expected := `TITLE: What is the meaning of life?
BODY: 42
METADATA: Don't panic
KEYWORDS: life, meaning, universe

TITLE: What is the answer to the ultimate question of life, the universe, and everything?
BODY: 42
METADATA: Don't panic
KEYWORDS: life, meaning, universe`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)	
}

func TestCardFromYaml(t *testing.T) {
	yamlString := `- title: What is the meaning of life?
  body: "42"
  metadata: Don't panic
  keywords:
  - life
  - meaning
  - universe
- title: What is the answer to the ultimate question of life, the universe, and everything?
  body: "42"
  metadata: Don't panic
  keywords:
  - life
  - meaning
  - universe
`

	cardList, err := data.CardListFromYaml(yamlString)

	if err != nil {
		t.Fatal(err)
	}

	if len(cardList) != 2 {
		t.Fatalf("unexpected length: %d", len(cardList))
	}

	fmt.Println("✅", len(cardList))
}

func TestCardListToYaml(t *testing.T) {
	cardList := []data.Card{
		{
			Title: "What is the meaning of life?",
			Body:  "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
		{
			Title: "What is the answer to the ultimate question of life, the universe, and everything?",
			Body:  "42",
			MetaData: "Don't panic",
			KeyWords: []string{"life", "meaning", "universe"},
		},
	}

	res := data.CardListToYaml(cardList)

	expected := `- title: What is the meaning of life?
  body: "42"
  metadata: Don't panic
  keywords:
  - life
  - meaning
  - universe
- title: What is the answer to the ultimate question of life, the universe, and everything?
  body: "42"
  metadata: Don't panic
  keywords:
  - life
  - meaning
  - universe
`

	if res != expected {
		t.Fatalf("unexpected result: %s", res)
	}

	fmt.Println("✅", res)

}
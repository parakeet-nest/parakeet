package content

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetArrayOfContentFiles(t *testing.T) {

	content, err := GetArrayOfContentFiles("./contents-for-test", ".txt")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("📝 content", content)

	if content[0] != "hello world" {
		t.Fatal("hello world not found")
	}
	if content[1] != "hey people" {
		t.Fatal("hey people not found")
	}

	if content[2] != "hello world" {
		t.Fatal("hello world not found")
	}
	if content[3] != "hey people" {
		t.Fatal("hey people not found")
	}

}

func TestGetMapOfContentFiles(t *testing.T) {

	content, err := GetMapOfContentFiles("./contents-for-test", ".txt")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("📝 content", content)

	if content["contents-for-test/01/hello.txt"] != "hello world" {
		t.Fatal("hello world not found")
	}
	if content["contents-for-test/02/hey.txt"] != "hey people" {
		t.Fatal("hey people not found")
	}
	if content["contents-for-test/01/hello.txt"] != "hello world" {
		t.Fatal("hello world not found")
	}
	if content["contents-for-test/02/hey.txt"] != "hey people" {
		t.Fatal("hey people not found")
	}

}

func TestGenerateContextFromDocs(t *testing.T) {
	content, err := GetArrayOfContentFiles("./contents-for-test", ".txt")

	if err != nil {
		t.Fatal(err)
	}

	context := GenerateContextFromDocs(content)

	fmt.Println("📝 context", context)

	if strings.Contains(context, "<doc>hello world</doc>") == false {
		t.Fatal("hello world not found")
	}
	if strings.Contains(context, "<doc>hey people</doc>") == false {
		t.Fatal("hey people not found")
	}

}


func TestInterpolation(t *testing.T) {

	human := struct {
		FirstName string
		LastName string

	}{
		FirstName: "Bob",
		LastName: "Morane",
	}

	tpl := `Hello I am {{.FirstName}} {{.LastName}}`

	res, err := InterpolateString(tpl, human)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("📙", res)

}

func TestSimpleChunker(t *testing.T) {

	rulesContent, err := ReadTextFile("./contents-for-test/rules/chronicles.md")
	if err != nil {
		t.Fatal(err)
	}
	chunks := ChunkText(rulesContent, 100, 15)

	for _, chunk := range chunks {
		fmt.Println("------------------------------------")
		fmt.Println(chunk)
		fmt.Println()
	}

}

func TestRegexSplitter(t *testing.T) {

	rulesContent, err := ReadTextFile("./contents-for-test/rules/chronicles.md")
	if err != nil {
		t.Fatal(err)
	}
	//chunks := SplitTextWithRegex(rulesContent, `#|\d+\.`)
	chunks := SplitTextWithRegex(rulesContent, `## *`)


	for _, chunk := range chunks {
		fmt.Println("=============================================")
		fmt.Println(chunk)
		fmt.Println()
	}

}



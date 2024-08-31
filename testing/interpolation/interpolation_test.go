package testing_interpolation

import (
	"fmt"
	"testing"
	"github.com/parakeet-nest/parakeet/content"
)

func TestInterpolation(t *testing.T) {

	human := struct {
		FirstName string
		LastName  string
	}{
		FirstName: "Bob",
		LastName:  "Morane",
	}

	tpl := `Hello I am {{.FirstName}} {{.LastName}}`

	res, err := content.InterpolateString(tpl, human)

	if err != nil {
		t.Fatal(err)
	}

	if res != "Hello I am Bob Morane" {
		t.Fatal("unexpected result")
	}

	fmt.Println("📙", res)

}

func TestInterpolationWithHTMLComments(t *testing.T) {
	chunk := content.Chunk{
		Prefix:       "##",
		Header:       "Title 2",
		Content:      "This is a content with <!-- a comment -->",
		Level:        2,
		ParentPrefix: "#",
		ParentLevel:  1,
		ParentHeader: "Title 1",
		Lineage:      "Title 1 > Title 2",
	}

	tpl :=
		`{{.Prefix}} {{.Header}}
<!-- Parent Section: {{.ParentPrefix}} {{.ParentHeader}} -->
<!-- Lineage: {{.Lineage}} -->
{{.Content}}`

	res, err := content.InterpolateString(tpl, chunk)

	if err != nil {
		t.Fatal(err)
	}

	expected := `## Title 2
<!-- Parent Section: # Title 1 -->
<!-- Lineage: Title 1 > Title 2 -->
This is a content with <!-- a comment -->`

	if res != expected {
		t.Fatal("unexpected result")
	}

	fmt.Println("📘", res)
}

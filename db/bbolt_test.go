package bbolt

import (
	"encoding/json"
	"fmt"
	"testing"
)

type Human struct {
	Name string  `json:"name"`
	Age  float64 `json:"age"`
}

func TestInitialize(t *testing.T) {
	fmt.Println("01-TestInitialize")
	db, err := Initialize("test.db", "test")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("🟢", db)
	defer db.Close()
}

func TestSaveAsJson(t *testing.T) {
	fmt.Println("02-TestSaveAsJson")
	db, err := Initialize("test.db", "test")
	if err != nil {
		t.Fatal(err)
	}
	bob := map[string]interface{}{
		"name": "Bob",
		"age":  42.0,
	}
	res, err := SaveAsJson(db, "test", "bob", bob)
	if err != nil {
		t.Fatal(err)
	}

	if res["name"] != "Bob" {
		t.Fatal("name is not Bob")
	}
	if res["age"].(float64) != 42 {
		t.Fatal("age is not 42")
	}
	defer db.Close()
	fmt.Println("🟢", res)
}

func TestGetFromJson(t *testing.T) {
	fmt.Println("03-TestGetFromJson")
	db, err := Initialize("test.db", "test")
	if err != nil {
		t.Fatal(err)
	}

	res, err := GetFromJson(db, "test", "bob")
	if err != nil {
		t.Fatal(err)
	}

	if res["name"] != "Bob" {
		t.Fatal("name is not Bob")
	}
	if res["age"].(float64) != 42 {
		t.Fatal("age is not 42")
	}
	defer db.Close()
	fmt.Println("🟢", res)
}

func TestGetFromJsonToHuman(t *testing.T) {
	fmt.Println("04-TestGetFromJsonToHuman")
	db, err := Initialize("test.db", "test")
	if err != nil {
		t.Fatal(err)
	}
	jsonStr := Get(db, "test", "bob")
	bob := Human{}
	err = json.Unmarshal([]byte(jsonStr), &bob)
	if err != nil {
		t.Fatal(err)
	}

	if bob.Name != "Bob" {
		t.Fatal("name is not Bob")
	}
	if bob.Age != 42 {
		t.Fatal("age is not 42")
	}
	defer db.Close()
	fmt.Println("🟢", bob)
}

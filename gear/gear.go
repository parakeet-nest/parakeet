package gear

import (
	"bytes"
	"encoding/json"
)

// PrettyString formats a JSON string into a more human-readable format.
//
// It takes a JSON string as input and returns a formatted JSON string and an error, if any.
// The formatted JSON string is indented with two spaces.
func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "  "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

// JSONParse parses a JSON string into a map[string]interface{}.
//
// It takes a JSON string as input and returns a map[string]interface{} and an error, if any.
// The map represents the parsed JSON object.
// If the JSON string cannot be parsed, the error is returned.
func JSONParse(str string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func JSONStringify(obj map[string]interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(b)
}
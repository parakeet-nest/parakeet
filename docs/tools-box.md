# Tools Box

## Strings and JSON

- `PrettyString` formats a JSON string into a more human-readable format: `gear.PrettyString(str string) (string, error)`
- `JSONParse` parses a JSON string into a `map[string]interface{}`: `gear.JSONParse(str string) (map[string]interface{}, error)`
- `JSONStringify` converts a `map[string]interface{}` object into a JSON string: `gear.JSONStringify(obj map[string]interface{}) string`

## Get and cast environment variable value at the same time:

- `gear.GetEnvFloat(key string, defaultValue float64) float64 `
- `gear.GetEnvInt(key string, defaultValue int) int`
- `gear.GetEnvString(key string, defaultValue string) string`

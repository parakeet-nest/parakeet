# How to set the Options of the query

The best way to set the options of the query is to start from the default Ollama options and override the fields you want to change.

```golang
options := llm.DefaultOptions()
// override the default value
options.Temperature = 0.5
```

Or you can use the `SetOptions` helper. This helper will set the default values for the fields not defined in the map:

Define only the fields you want to override:
```golang
options := llm.SetOptions(map[string]interface{}{
  "Temperature": 0.5,
})
```

Or use the `SetOptions` helper with the `option` enums:
```golang
options := llm.SetOptions(map[string]interface{}{
  option.Temperature: 0.5,
  option.RepeatLastN: 2,
})
```

!!! note
    Before, the JSON serialization of the `Options` used the `omitempty` tag.

    The `omitempty` tag prevents a field from being serialised if its value is the zero value for the field's type (e.g., 0.0 for float64).

    That means when `Temperature` equals `0.0`, the field is not serialised (then Ollama will use the `Temperature` default value, which equals `0.8`).

    The problem will happen for every value equal to `0` or `0.0`

    Since now, the `omitempty` tag is removed from the `Options` struct.
# UI helpers

!!! info "📦 `ui` package"


These helpers provide methods to help you to create a better CLI interface for interacting with the LLMs.

## `Input`

`func Input(color, prompt string) (string, error)`

`Input` displays a prompt with the specified color and waits for user input.

**Parameters**:
  - color: A string representing the color of the prompt text.
  - prompt: A string representing the prompt message to display.

**Returns**:
  - A string containing the user input, trimmed of any leading or trailing whitespace.
  - An error if there was an issue running the input program or if the input could not be retrieved.

```go
res, err := ui.Input(colors.Cyan, "🤖 ask me something>")
```

## `Println`

`Println(color string, strs ...interface{})`

`Println` prints the provided strings with the specified color using the lipgloss styling library.
The color parameter should be a string representing the desired color.
The strs parameter is a variadic argument that accepts multiple values to be printed.

**Parameters**:
  - color: A string representing the color to be used for the text.
  - strs: A variadic parameter that accepts multiple values to be printed.

```go
ui.Println(colors.Magenta, "👋 hello 🤖")
```

Imports:
```go
"github.com/parakeet-nest/parakeet/ui"
"github.com/parakeet-nest/parakeet/ui/colors"
```

!!! info
    I used these two great libraries to create the helpers:
    - github.com/charmbracelet/bubbletea
    - github.com/charmbracelet/lipgloss

!!! note
	👀 you will find complete examples in:

    - [examples/58-michael-burnham](https://github.com/parakeet-nest/parakeet/tree/main/examples/58-michael-burnham)
    - [examples/59-jean-luc-picard-contextual-retrieval](https://github.com/parakeet-nest/parakeet/tree/main/examples/59-jean-luc-picard-contextual-retrieval)


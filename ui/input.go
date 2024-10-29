package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textInput textinput.Model
	err       error
}

// initialModel initializes and returns a new model with a text input field.
// The text input field is configured with a placeholder, focus, character limit, width, and a prompt.
//
// Parameters:
//   - prompt: A string that sets the prompt for the text input field.
//
// Returns:
//   - model: A struct containing the configured text input field and an error field.
func initialModel(prompt string) model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 255
	ti.Width = 80

	ti.Prompt = prompt

	return model{
		textInput: ti,
		err:       nil,
	}
}

// Init initializes the model and returns a command to start the text input blinking.
// It is part of the tea.Model interface implementation.
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles incoming messages and updates the model accordingly.
// It processes key messages such as Enter, Ctrl+C, and Esc to quit the application.
// It also handles errors by storing them in the model.
// Finally, it updates the text input component and returns the updated model and command.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// Handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}


var promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#008000"))
// View renders the current state of the text input model as a string.
// It applies the promptStyle to the rendered view of the text input and
// appends a newline character.
//
// Returns:
//   A string representing the styled view of the text input.
func (m model) View() string {
	return promptStyle.Render(m.textInput.View()+"\n") 
}


// Input displays a prompt with the specified color and waits for user input.
//
// Parameters:
//   - color: A string representing the color of the prompt text.
//   - prompt: A string representing the prompt message to display.
//
// Returns:
//   - A string containing the user input, trimmed of any leading or trailing whitespace.
//   - An error if there was an issue running the input program or if the input could not be retrieved.
func Input(color, prompt string) (string, error) {
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	p := tea.NewProgram(initialModel(prompt))
	m, err := p.Run()
	if err != nil {
		return "", err
	}
	if m, ok := m.(model); ok {
		return strings.TrimSpace(m.textInput.Value()), nil
	}
	return "", fmt.Errorf("😡 unable to get input")
}

// Println prints the provided strings with the specified color using the lipgloss styling library.
// The color parameter should be a string representing the desired color.
// The strs parameter is a variadic argument that accepts multiple values to be printed.
//
// Parameters:
//   - color: A string representing the color to be used for the text.
//   - strs: A variadic parameter that accepts multiple values to be printed.
//
// Example usage:
//   Println(colors.Magenta, "Hello", "world!")
func Println(color string, strs ...interface{}) {
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	
	// Convert all arguments to strings
	strSlice := make([]string, len(strs))
	for i, v := range strs {
		strSlice[i] = fmt.Sprint(v)
	}
	
	// Join all strings and render with the style
	renderedString := textStyle.Render(strings.Join(strSlice, " "))
	
	// Print the rendered string with a newline
	fmt.Println(renderedString)
}

func Print(color string, strs ...interface{}) {
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	
	// Convert all arguments to strings
	strSlice := make([]string, len(strs))
	for i, v := range strs {
		strSlice[i] = fmt.Sprint(v)
	}
	
	// Join all strings and render with the style
	renderedString := textStyle.Render(strings.Join(strSlice, " "))
	
	// Print the rendered string with a newline
	fmt.Print(renderedString)
}

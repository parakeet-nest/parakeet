package cli

import (
	"fmt"
	"os"
	"strings"
)

type Arg struct {
	Name string
}

type Flag struct {
	Name  string
	Value string
}

// Settings parses command-line arguments and flags.
//
// It skips the program name and processes the remaining arguments.
// Arguments that start with "--" are considered flags, and the function
// checks if the next argument is a value for the flag. If so, it pairs
// the flag with its value; otherwise, it pairs the flag with an empty string.
// Arguments that do not start with "--" are considered positional arguments.
//
// Returns two slices: one containing the positional arguments and the other
// containing the flags with their respective values.
func Settings() ([]Arg, []Flag) {
	args := os.Args[1:] // Skip the program name
	var arguments []Arg
	var flags []Flag

	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			flagName := strings.TrimPrefix(args[i], "--")
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				flags = append(flags, Flag{Name: flagName, Value: args[i+1]})
				i++ // Skip the next argument as it's the flag value
			} else {
				flags = append(flags, Flag{Name: flagName, Value: ""})
			}
		} else {
			arguments = append(arguments, Arg{Name: args[i]})
		}
	}

	return arguments, flags
}

// FlagValue retrieves the value of a flag by its name from a slice of Flag structs.
// If the flag is not found, it returns an empty string.
//
// Parameters:
//   - name: The name of the flag to search for.
//   - flags: A slice of Flag structs to search within.
//
// Returns:
//   The value of the flag if found, otherwise an empty string.
func FlagValue(name string, flags []Flag) string {
	for _, flag := range flags {
		if flag.Name == name {
			return flag.Value
		}
	}
	return ""
}

// HasArg checks if an argument with the specified name exists in the provided slice of arguments.
// 
// Parameters:
// - name: The name of the argument to search for.
// - args: A slice of Arg structures to search within.
//
// Returns:
// - bool: True if an argument with the specified name is found, otherwise false.
func HasArg(name string, args []Arg) bool {
	for _, arg := range args {
		if arg.Name == name {
			return true
		}
	}
	return false
}

// HasFlag checks if a flag with the specified name exists in the provided slice of flags.
// 
// Parameters:
// - name: The name of the flag to search for.
// - flags: A slice of Flag objects to search within.
//
// Returns:
// - bool: True if a flag with the specified name is found, otherwise false.
func HasFlag(name string, flags []Flag) bool {
	for _, flag := range flags {
		if flag.Name == name {
			return true
		}
	}
	return false
}

// ArgsTail extracts the names from a slice of Arg structs and returns them as a slice of strings.
// 
// Parameters:
// - args: A slice of Arg structs from which the names will be extracted.
//
// Returns:
// - A slice of strings containing the names of the provided Arg structs.
func ArgsTail(args []Arg) []string {
	names := make([]string, len(args))
	for i, arg := range args {
		names[i] = arg.Name
	}
	return names
}

// FlagsTail takes a slice of Flag structs and returns a slice of strings
// containing the names of those flags.
//
// Parameters:
//   flags []Flag: A slice of Flag structs.
//
// Returns:
//   []string: A slice of strings containing the names of the flags.
func FlagsTail(flags []Flag) []string {
	names := make([]string, len(flags))
	for i, flag := range flags {
		names[i] = flag.Name
	}
	return names
}

// FlagsWithNamesTail takes a slice of Flag structs and returns a slice of strings,
// where each string is a formatted pair of the flag's name and value in the form "name=value".
//
// Parameters:
//   flags []Flag - A slice of Flag structs, each containing a Name and a Value.
//
// Returns:
//   []string - A slice of strings, each representing a flag's name and value pair.
func FlagsWithNamesTail(flags []Flag) []string {
	pairs := make([]string, len(flags))
	for i, flag := range flags {
		pairs[i] = fmt.Sprintf("%s=%s", flag.Name, flag.Value)
	}
	return pairs
}

// containsSubsequence checks if the subSeq slice is a subsequence of the mainSeq slice.
// A subsequence is a sequence that appears in the same relative order, but not necessarily consecutively.
// 
// Parameters:
// - mainSeq: The main sequence of strings to be checked.
// - subSeq: The subsequence of strings to be searched for within the main sequence.
//
// Returns:
// - bool: true if subSeq is a subsequence of mainSeq, false otherwise.
func containsSubsequence(mainSeq, subSeq []string) bool {
	if len(subSeq) == 0 {
		return true
	}
	if len(mainSeq) == 0 {
		return false
	}

	subIdx := 0
	for _, item := range mainSeq {
		if item == subSeq[subIdx] {
			subIdx++
			if subIdx == len(subSeq) {
				return true
			}
		}
	}
	return false
}

// HasSubsequence checks if the given subsequence of strings (subSeq) is present
// in the tail of the provided arguments (args).
//
// Parameters:
//   - args: A slice of Arg representing the arguments to be checked.
//   - subSeq: A slice of strings representing the subsequence to look for.
//
// Returns:
//   - bool: True if the subsequence is found in the tail of the arguments, false otherwise.
func HasSubsequence(args []Arg, subSeq []string) bool {
	return containsSubsequence(ArgsTail(args), subSeq)
}
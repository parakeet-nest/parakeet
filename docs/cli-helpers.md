# CLI Helpers

!!! info "📦 `content` package"

These helpers provide methods to parse and use command-line arguments and flags.

## `Settings`

`Settings` parses command-line arguments and flags.

It skips the program name and processes the remaining arguments.
Arguments that start with "--" are considered flags, and the function
checks if the next argument is a value for the flag. If so, it pairs
the flag with its value; otherwise, it pairs the flag with an empty string.
Arguments that do not start with "--" are considered positional arguments.

**Returns** two slices: one containing the positional arguments and the other
containing the flags with their respective values.

**Example**:

```go
// default values
ollamaUrl := "http://localhost:11434"
chatModel := "llama3.1:8b"
embeddingsModel := "bge-m3:latest"

args, flags := cli.Settings()

if cli.HasFlag("url", flags) {
    ollamaUrl = cli.FlagValue("url", flags)
}

if cli.HasFlag("chat-model", flags) {
    chatModel = cli.FlagValue("chat-model", flags)
}

if cli.HasFlag("embeddings-model", flags) {
    embeddingsModel = cli.FlagValue("embeddings-model", flags)
}

switch cmd := cli.ArgsTail(args); cmd[0] {
case "create-embeddings":
    fmt.Println(embeddingsModel)
case "chat":
    fmt.Println(chatModel)
default:
    fmt.Println("Unknown command:", cmd[0])
}
```

## `FlagValue`

`FlagValue` retrieves the value of a flag by its name from a slice of Flag structs.
If the flag is not found, it returns an empty string.

**Parameters**:
  - `name`: The name of the flag to search for.
  - `flags`: A slice of Flag structs to search within.

**Returns**:
  The value of the flag if found, otherwise an empty string.

## `HasArg`

`HasArg` checks if an argument with the specified name exists in the provided slice of arguments.

**Parameters**:
- `name`: The name of the argument to search for.
- `args`: A slice of Arg structures to search within.

**Returns**:
- `bool`: True if an argument with the specified name is found, otherwise false.

## `HasFlag`

`HasFlag` checks if a flag with the specified name exists in the provided slice of flags.

**Parameters**:
- `name`: The name of the flag to search for.
- `flags`: A slice of Flag objects to search within.

**Returns**:
- `bool`: True if a flag with the specified name is found, otherwise false.

## `ArgsTail`

`ArgsTail` extracts the names from a slice of Arg structs and returns them as a slice of strings.

**Parameters**:
- `args`: A slice of Arg structs from which the names will be extracted.

**Returns**:
- A slice of strings containing the names of the provided Arg structs.

## `FlagsTail`

`FlagsTail` takes a slice of Flag structs and returns a slice of strings containing the names of those flags.

**Parameters**:
  `flags []Flag`: A slice of Flag structs.

**Returns**:
  `[]string`: A slice of strings containing the names of the flags.

## `FlagsWithNamesTail`

`FlagsWithNamesTail` takes a slice of Flag structs and returns a slice of strings, where each string is a formatted pair of the flag's name and value in the form "name=value".

**Parameters**:
  `flags []Flag` - A slice of Flag structs, each containing a Name and a Value.

**Returns**:
  `[]string` - A slice of strings, each representing a flag's name and value pair.

## `HasSubsequence`

`HasSubsequence` checks if the given subsequence of strings (subSeq) is present in the tail of the provided arguments (args).

Parameters:
  - `args`: A slice of Arg representing the arguments to be checked.
  - `subSeq`: A slice of strings representing the subsequence to look for.

Returns:
  - `bool`: True if the subsequence is found in the tail of the arguments, false otherwise.



!!! note
	👀 you will find complete examples in:

    - [examples/59-jean-luc-picard-contextual-retrieval](https://github.com/parakeet-nest/parakeet/tree/main/examples/59-jean-luc-picard-contextual-retrieval)


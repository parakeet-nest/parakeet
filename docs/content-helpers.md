# Content Helpers

!!! info "🚧 work in progress"

## Context construction

- `GenerateContextFromDocs` generates the context content from a slice of documents: `content.GenerateContextFromDocs(docs []string) string`

## Read and Write files

- `content.ReadTextFile(path string) (string, error)`
- `content.WriteTextFile(path, content string) error`

## Parsing path

- `FindFiles` searches for files with a specific extension in the given root directory and its subdirectories: `content.FindFiles(dirPath string, ext string) ([]string, error)`
  - Returns:
    - []string: A slice of file paths that match the given extension.
    - error: An error if the search encounters any issues.
- `ForEachFile` iterates over all files with a specific extension in a directory and its subdirectories: `content.ForEachFile(dirPath string, ext string, callback func(string) error) ([]string, error)`
  - Returns:
    - []string: A slice of file paths that match the given extension.
    - error: An error if the search encounters any issues.
- `GetArrayOfContentFiles` searches for files with a specific extension in the given directory and its subdirectories: `content.GetArrayOfContentFiles(dirPath string, ext string) ([]string, error)`
- `GetMapOfContentFiles` searches for files with a specific extension in the given directory and its subdirectories: `content.GetMapOfContentFiles(dirPath string, ext string) (map[string]string, error)`


## String interpolation

- `content.InterpolateString(str string, vars interface{}) (string, error)`

!!! note
	👀 you will find a complete example in:

    - [examples/40-rag-with-elastic-markdown](https://github.com/parakeet-nest/parakeet/tree/main/examples/40-rag-with-elastic-markdown)
    - [examples/46-create-an-expert](https://github.com/parakeet-nest/parakeet/tree/main/examples/46-create-an-expert)


## Estimate the number of tokens in a text

- `content.CountTokens(text string) int`
- `content.CountTokensAdvanced(text string) int`
- `content.EstimateGPTTokens(text string) int`

> this could be useful to estimate the value of `num_ctx`

!!! note
	👀 you will find a complete example in:

    - [examples/72-gitingest-es](https://github.com/parakeet-nest/parakeet/tree/main/examples/72-gitingest-es)
    - [examples/73-gitingest-daphnia](https://github.com/parakeet-nest/parakeet/tree/main/examples/73-gitingest-daphnia)
    - [examples/74-rag-with-daphnia](https://github.com/parakeet-nest/parakeet/tree/main/examples/74-rag-with-daphnia)
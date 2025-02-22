# Parsing Source Code

!!! info "🚧 work in progress"

## Extract elements from source code:

- `source.ExtractCodeElements(fileContent string, language string) ([]CodeElement, error)`

```golang
// CodeElement represents a code structure element (class, function, method)
type CodeElement struct {
	Type        string // "class", "function", "method"
	Name        string
	Signature   string
	Description string
	LineNumber  int
	ParentClass string // For methods
	Parameters  []string
	Source      string // Source code of the element
}
```

> the `Signature` could be useful to add context to embeddings.


!!! note
	👀 you will find a complete example in:

    - [examples/72-gitingest-es](https://github.com/parakeet-nest/parakeet/tree/main/examples/72-gitingest-es)
    - [examples/73-gitingest-daphnia](https://github.com/parakeet-nest/parakeet/tree/main/examples/73-gitingest-daphnia)
    - [examples/74-rag-with-daphnia](https://github.com/parakeet-nest/parakeet/tree/main/examples/74-rag-with-daphnia)


    
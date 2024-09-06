# Chunkers and Splitters

There are several methods in the `content` package to help you chunk and split text:

- `ChunkText` takes a text string and divides it into chunks of a specified size with a given overlap. It returns a slice of strings, where each string represents a chunk of the original text.

```golang
chunks := content.ChunkText(documentContent, 900, 400)
```

- `SplitTextWithDelimiter` splits the given text using the specified delimiter and returns a slice of strings.

```golang
chunks := content.SplitTextWithDelimiter(documentContent, "<!-- SPLIT -->")
```

- `SplitTextWithRegex` splits the given text using the provided regular expression delimiter. It returns a slice of strings containing the split parts of the text.

```golang
chunks := content.SplitTextWithRegex(documentContent, `## *`)
```

- `SplitMarkdownBySections` splits the given markdown text using the title sections (`#, ##, etc.`) and returns a slice of strings.

```golang
chunks := content.SplitMarkdownBySections(documentContent)
```

- `SplitAsciiDocBySections` splits the given asciidoc text using the title sections (`=, ==, etc.`) and returns a slice of strings.

```golang
chunks := content.SplitAsciiDocBySections(documentContent)
```

- `SplitHTMLBySections` splits the given html text using the title sections (`h1, h2, h3, h4, h5, h6`) and returns a slice of strings.

```golang
chunks := content.SplitHTMLBySections(documentContent)
```
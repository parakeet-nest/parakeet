# Parsing Markdown

!!! info "🚧 work in progress" (TODO: add more examples)

I created some helpers to help to split markdown documents and create better chunks.

- `ParseMarkdownWithHierarchy` chunks a markdown document while maintaining semantic meaning and preserving the relationship between sections.

```golang
chunks := content.ParseMarkdownWithHierarchy(document)
```
> `func ParseMarkdownWithHierarchy(document string) []Chunk`

You will get the following data:
```golang
chunk := Chunk{
    Level:        level,
    Prefix:       prefix,
    Header:       header,
    Content:      strings.TrimSpace(content),
    ParentPrefix: parent.Prefix,
    ParentLevel:  parent.Level,
    ParentHeader: parent.Header,
}
```
Then you can add meta data when creating the vectors thanks to these fields: `ParentPrefix`, `ParentLevel`, `ParentHeader`.

- `ParseMarkdownWithLineage` parses the given markdown content and returns a slice of Chunk structs. Each Chunk represents a header and its associated content, along with its **hierarchical lineage**.

```golang
chunks := content.ParseMarkdownWithLineage(document)
```
> `func ParseMarkdownWithLineage(document string) []Chunk`

You will get the following data:
```golang
chunk := Chunk{
    Level:        level,
    Prefix:       prefix,
    Header:       header,
    Content:      strings.TrimSpace(content),
    ParentPrefix: parent.Prefix,
    ParentLevel:  parent.Level,
    ParentHeader: parent.Header,
    Lineage:      lineage,
}
```
Then you can add meta data when creating the vectors thanks to this field: `Lineage`.

`Lineage` will keep the path of the sections. For example, with this document:

```markdown
# Tiefling Species in Fantasy Realms: A Comprehensive Analysis

... some text ...

## Professional Development and Education

... some text ...
```

The `Lineage` value of the chunk of the second section will be:

```raw
Tiefling Species in Fantasy Realms: A Comprehensive Analysis > Professional Development and Education
```

!!! note
	👀 you will find a complete example in:

    - [examples/65-hyde](https://github.com/parakeet-nest/parakeet/tree/main/examples/65-hyde)
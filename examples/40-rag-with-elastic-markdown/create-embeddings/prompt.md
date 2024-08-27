

This is my markdown content:

```markdown
# This is level 1 section

bla bla bla

## This is level 2 section

bla bla bla
bla bla bla
bla bla bla

### This is level 3 section

bla bla bla
bla bla bla
bla bla bla

#### This is level 4 section

bla bla bla
bla bla bla
bla bla bla

## This is another level 2 section

bla bla bla
bla bla bla

### This is another level 3 section

bla bla bla
bla bla bla
bla bla bla

#### This is another level 4 section

bla bla bla
bla bla bla
bla bla bla

#### This is another level 4 section

bla bla bla
bla bla bla
bla bla bla

##### This is another level 5 section

bla bla bla
bla bla bla
bla bla bla

```

I want to split the markdown content into sections and subsections.
I use Golang.

Every result (chunk) will be a struct with the following fields:

```golang
type Chunk struct {
	Level        int
	Prefix       string
	Header       string
	Content      string
    ParentPrefix string
    ParentLevel  int
    ParentHeader string
}
```


For example, the chunk for the first section will be:

```golang
Chunk{
    Level:   1,
    Prefix:  "#",
    Header:  "This is level 1 section",
    Content: "bla bla bla",
    ParentPrefix: "",
    ParentLevel:  0,
    ParentHeader: "",
}
```

The chunk for the second section will be:

```golang
Chunk{
    Level:   2,
    Prefix:  "##",
    Header:  "This is level 2 section",
    Content: "bla bla bla\nbla bla bla\nbla bla bla",
    ParentPrefix: "#",
    ParentLevel:  1,
    ParentHeader: "This is level 1 section",
}
```

The chunk for the third section will be:

```golang
Chunk{
    Level:   3,
    Prefix:  "###",
    Header:  "This is level 3 section",
    Content: "bla bla bla\nbla bla bla\nbla bla bla",
    ParentPrefix: "##",
    ParentLevel:  2,
    ParentHeader: "This is level 2 section",
}
```

The chunk for the fourth section will be:

```golang
Chunk{
    Level:   4,
    Prefix:  "####",
    Header:  "This is level 4 section",
    Content: "bla bla bla\nbla bla bla\nbla bla bla",
    ParentPrefix: "###",
    ParentLevel:  3,
    ParentHeader: "This is level 3 section",
}
```

The chunk for the another level 2 section will be:

```golang
Chunk{
    Level:   2,
    Prefix:  "##",
    Header:  "This is another level 2 section",
    Content: "bla bla bla\nbla bla bla",
    ParentPrefix: "#",
    ParentLevel:  1,
    ParentHeader: "This is level 1 section",
}
```

The chunk for the another level 3 section will be:

```golang
Chunk{
    Level:   3,
    Prefix:  "###",
    Header:  "This is another level 3 section",
    Content: "bla bla bla\nbla bla bla\nbla bla bla",
    ParentPrefix: "##",
    ParentLevel:  2,
    ParentHeader: "This is another level 2 section",
}
```

The chunk for the another level 4 section will be:

```golang
Chunk{
    Level:   4,
    Prefix:  "####",
    Header:  "This is another level 4 section",
    Content: "bla bla bla\nbla bla bla\nbla bla bla",
    ParentPrefix: "###",
    ParentLevel:  3,
    ParentHeader: "This is another level 3 section",
}
```

The chunk for the another level 5 section will be:

```golang
Chunk{
    Level:   5,
    Prefix:  "#####",
    Header:  "This is another level 5 section",
    Content: "bla bla bla\nbla bla bla\nbla bla bla",
    ParentPrefix: "####",
    ParentLevel:  4,
    ParentHeader: "This is another level 4 section",
}
```

Generate the function `SplitMarkdown` that receives a markdown content and returns a slice of chunks like the above examples.



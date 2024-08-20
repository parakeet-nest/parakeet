# RAG with the README file
> 🚧 wip

## First: add the metadata and splitters

I added (manually) metadata to the `README.md` file a the begining of every important section.
I used Llama3 to generate the metadata with this system instructions:

```text
You are a programming expert. When analysing a given Markdown document, you are able:
- To determine the main topic the document
- To make a brief summary of the document
- To extract the most accurate key words from the document
The output will use the following format:
<!--
TOPIC: <the topic of the piece of the document>
SUMMARY: <the summary of the piece of the document>
KEYWORDS: <the accurate keywords from the piece of the document>
-->
Analyse the user document and extract the informations according the instructions with the given output format:
```

I added (manually) splitters (`<!-- split -->`) to the `README.md` file a the end of every important section.


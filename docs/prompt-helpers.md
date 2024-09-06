
<!-- TOPIC: Prompt helpers and meta prompts SUMMARY: A collection of special instructions, known as meta-prompts, to guide language models in generating specific kinds of responses. KEYWORDS: Meta prompts, prompt helpers, AI, LLM, natural language processing, NLP -->
# Prompt helpers

## Meta prompts
> package: `prompt`

Meta-prompts are special instructions embedded within a prompt to guide a language model in generating a specific kind of response.

|  Meta-Prompt   |  Purpose  |
| :------------  | :-------- |
|[Brief] What is AI? | For a concise answer
|[In Layman’s Terms] Explain LLM | For a simplified explanation
|[As a Story] Describe the evolution of cars | To get the information in story form
|[Pros and Cons] Is AI useful? | For a balanced view with advantages and disadvantages
|[Step-by-Step] How to do a smart prompt? | For a detailed, step-by-step guide
|[Factual] What is the best pizza of the world? | For a straightforward, factual answer
|[Opinion] What is the best pizza of the world? | To get an opinion-based answer
|[Comparison] Compare pineapple pizza to pepperoni pizza | For a comparative analysis
|[Timeline] What are the key milestones to develop a WebApp? | For a chronological account of key events
|[As a Poem] How to cook a cake? | For a poetic description
|[For Kids] How to cook a cake? | For a child-friendly explanation
|[Advantages Only] What are the benefits of AI? | To get a list of only the advantages
|[As a Recipe] How to cook a cake? | To receive the information in the form of a recipe

## Meta prompts methods

- `prompt.Brief(s string) string`
- `prompt.InLaymansTerms(s string) string`
- `prompt.AsAStory(s string) string`
- `prompt.ProsAndCons(s string) string`
- `prompt.StepByStep(s string) string`
- `prompt.Factual(s string) string`
- `prompt.Opinion(s string) string`
- `prompt.Comparison(s string) string`
- `prompt.Timeline(s string) string`
- `prompt.AsAPoem(s string) string`
- `prompt.ForKids(s string) string`
- `prompt.AdvantagesOnly(s string) string`
- `prompt.AsARecipe(s string) string`


package prompt

/*
#### Meta prompts
Meta-prompts are special instructions embedded within a prompt to guide a language model in generating a specific kind of response.

|  Meta-Prompt   |  Purpose  |
| :------------  | :-------- |
|[Brief] What is Kubernetes? | For a concise answer
|[In Layman’s Terms] Explain machine learning | For a simplified explanation
|[As a Story] Describe the evolution of programming languages | To get the information in story form
|[Pros and Cons] Should our company move to the cloud? | For a balanced view with advantages and disadvantages
|[Step-by-Step] How do I set up a VPN? | For a detailed, step-by-step guide
|[Factual] What is the current version of Python? | For a straightforward, factual answer
|[Opinion] Which is better for backend development: Node.js or Django? | To get an opinion-based answer
|[Comparison] Compare SQL databases and NoSQL databases| For a comparative analysis
|[Timeline] What are the key milestones in cybersecurity? | For a chronological account of key events
|[As a Poem] Tell me about coding | For a poetic description
|[For Kids] How does the internet work? | For a child-friendly explanation
|[Advantages Only] What are the benefits of using containers? | To get a list of only the advantages
|[As a Recipe] How to write a Python script for automating tasks? | To receive the information in the form of a recipe
*/

func Brief(s string) string {
	return "[Brief] " + s
}

func InLaymansTerms(s string) string {
	return "[In Layman’s Terms] " + s
}

func AsAStory(s string) string {
	return "[As a Story] " + s
}

func ProsAndCons(s string) string {
	return "[Pros and Cons] " + s
}

func StepByStep(s string) string {
	return "[Step-by-Step] " + s
}

func Factual(s string) string {
	return "[Factual] " + s
}

func Opinion(s string) string {
	return "[Opinion] " + s
}

func Comparison(s string) string {
	return "[Comparison] " + s
}

func Timeline(s string) string {
	return "[Timeline] " + s
}

func AsAPoem(s string) string {
	return "[As a Poem] " + s
}

func ForKids(s string) string {
	return "[For Kids] " + s
}

func AdvantagesOnly(s string) string {
	return "[Advantages Only] " + s
}

func AsARecipe(s string) string {
	return "[As a Recipe] " + s
}

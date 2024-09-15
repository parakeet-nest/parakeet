package main

import (
	"github.com/parakeet-nest/parakeet/completion"
	"github.com/parakeet-nest/parakeet/content"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/enums/option"

	"fmt"
	"log"
)

/*
https://github.com/01-ai/Yi-Coder/blob/main/cookbook/System_prompt/System_prompt.ipynb
*/

// This scenario demonstrates how Yi-Coder can identify errors and insert the correct code to fix them.

func main() {
	ollamaUrl := "http://localhost:11434"
	model := "yi-coder:9b"
	//model := "yi-coder:1.5b"

	systemContent := `SYSTEM:
	You are Yi-Coder, you are exceptionally skilled in programming, coding, devops and any computer-related issues.
	`

	allSourceCode, err := content.GetMapOfContentFiles("./", ".yml")
	if err != nil {
		log.Fatal(err)
	}

	codebase := "CODEBASE:\n"
	for _, pipeline := range allSourceCode {
		codebase += "<>\n```yaml\n" + pipeline + "```\n<>\n"
	}

	//userContent := `[Step-by-Step] Using the above codebase, explain the GitLab CI pipeline.
	//Make a detqiled section per CI job.`

	userContent := `Using the above the GitLab CI pipeline configuration in YAML format. 
	Please analyze this pipeline step by step and explain the following details:

      1. **Overview of the Pipeline**: Provide a high-level summary of the purpose and structure of the pipeline.
      2. **Pipeline Stages**: List and describe each stage defined in the pipeline, explaining their purpose and order of execution.
      3. **Jobs within Stages**: For each stage, detail the jobs that are included, explaining the commands or scripts they execute, and how they contribute to the overall pipeline.
      4. **Triggers and Conditions**: Identify any triggers, conditions, or rules that dictate when the pipeline runs or when specific jobs are executed.
      5. **Artifacts and Dependencies**: Explain how artifacts are used or shared between jobs, and describe any dependencies between jobs or stages.
      6. **Environment Configuration**: Detail any environment-specific configurations, variables, or secrets used within the pipeline.
      7. **Error Handling and Notifications**: Describe how the pipeline handles errors and if there are any notifications set up for failed jobs or stages.
      8. **Optimization or Improvements**: Suggest potential optimizations or improvements to enhance the efficiency or reliability of the pipeline.

    Please provide the explanation in a clear, concise manner, with examples where applicable.`

	options := llm.SetOptions(map[string]interface{}{
		option.Temperature: 0.0,
	})

	query := llm.GenQuery{
		Model:   model,
		Prompt:  systemContent + codebase + userContent,
		Options: options,
	}

	_, err = completion.GenerateStream(ollamaUrl, query,
		func(answer llm.GenAnswer) error {
			fmt.Print(answer.Response)
			return nil
		})

	if err != nil {
		log.Fatal("😡:", err)
	}

}

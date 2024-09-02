# This is a WIP 🚧

https://ollama.com/blog/tool-support

TODO: update the API Call to use the "tools" parameter instead of "system"
https://github.com/ollama/ollama/blob/main/docs/api.md#chat-request-with-tools

Add `Tools` field to llm.Query{}


---
https://huggingface.co/datasets/Salesforce/xlam-function-calling-60k


👋 use "tools" instead of "system" (is it known by ollama?)

format:
```json
{
  "query": "Find the sum of all the multiples of 3 and 5 between 1 and 1000. Also find the product of the first five prime numbers.",
  "tools": [
    {
      "name": "math_toolkit.sum_of_multiples",
      "description": "Find the sum of all multiples of specified numbers within a specified range.",
      "parameters": {
        "lower_limit": {
          "type": "int",
          "description": "The start of the range (inclusive).",
          "required": true
        },
        "upper_limit": {
          "type": "int",
          "description": "The end of the range (inclusive).",
          "required": true
        },
        "multiples": {
          "type": "list",
          "description": "The numbers to find multiples of.",
          "required": true
        }
      }
    },
    {
      "name": "math_toolkit.product_of_primes",
      "description": "Find the product of the first n prime numbers.",
      "parameters": {
        "count": {
          "type": "int",
          "description": "The number of prime numbers to multiply together.",
          "required": true
        }
      }
    }
  ],
  "answers": [
    {
      "name": "math_toolkit.sum_of_multiples",
      "arguments": {
        "lower_limit": 1,
        "upper_limit": 1000,
        "multiples": [3, 5]
      }
    },
    {
      "name": "math_toolkit.product_of_primes",
      "arguments": {
        "count": 5
      }
    }
  ]
}

```
package mcpsse

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/tools"
)

type Client struct {
	mcpClient *client.SSEMCPClient
	BaseURL   string
	ctx       context.Context
}

func NewClient(ctx context.Context, baseURL string, options ...string) (Client, error) {
	
	var bearerToken string
	if len(options) > 0 {
		bearerToken = options[0]
	} else {
		bearerToken = ""
	}
	var mcpClient *client.SSEMCPClient
	var err error
	if bearerToken == "" {
		mcpClient, err = client.NewSSEMCPClient(baseURL + "/sse")

	} else {

		mcpClient, err = client.NewSSEMCPClient(
			baseURL+"/sse",
			client.WithHeaders(map[string]string{
				"Authorization": "Bearer " + bearerToken}),
		)
	}

	if err != nil {
		return Client{}, &SSEClientCreationError{Message: fmt.Sprintf("Failed to create client: %v", err)}
	}
	defer mcpClient.Close()

	return Client{
		mcpClient: mcpClient,
		BaseURL:   baseURL,
		ctx:       ctx,
	}, nil
}

func (c *Client) Start() error {
	if err := c.mcpClient.Start(c.ctx); err != nil {
		return &SSEClientStartError{Message: fmt.Sprintf("Failed to start client: %v", err)}
	}
	return nil
}

func (c *Client) Initialize() (*mcp.InitializeResult, error) {
	initRequest := mcp.InitializeRequest{}

	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "parakeet mcp client",
		Version: "1.0.0",
	}

	initResult, err := c.mcpClient.Initialize(c.ctx, initRequest)
	if err != nil {
		return nil, &SSEClientInitializationError{Message: fmt.Sprintf("Failed to initialize client: %v", err)}
	}

	return initResult, nil
}

func (c *Client) ListTools() ([]llm.Tool, error) {
	toolsRequest := mcp.ListToolsRequest{}
	mcpTools, err := c.mcpClient.ListTools(c.ctx, toolsRequest)

	if err != nil {
		return nil, &SSEGetToolsError{Message: fmt.Sprintf("Failed to list tools: %v", err)}
	}

	// Convert mcp.Tools to llm.Tools
	// TODO: Handle errors during conversion
	ollamaTools := tools.ConvertMCPTools(mcpTools.Tools)
	return ollamaTools, nil
}

// string or interface?
type CallToolResult struct {
	Text string
	Type string
}

func (c *Client) CallTool(functionName string, arguments map[string]interface{}) (CallToolResult, error) {
	toolRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
	}
	toolRequest.Params.Name = functionName
	toolRequest.Params.Arguments = arguments

	mcpResult, err := c.mcpClient.CallTool(c.ctx, toolRequest)
	if err != nil {
		return CallToolResult{}, &SSEToolCallError{Message: fmt.Sprintf("Failed to call tool: %v", err)}
	}

	text, err := getTextFromResult(mcpResult)
	if err != nil {
		return CallToolResult{}, &SSEToolCallError{Message: fmt.Sprintf("Failed to extract text from result: %v", err)}
	}

	//! assumption: the result is always text
	return CallToolResult{Text: text, Type: "text"}, nil
}

func (c *Client) Close() error {
	return c.mcpClient.Close()
}

// getTextFromResult extracts the text content from an mcp.CallToolResult object.
// It returns the extracted text as a string and an error if the extraction fails.
//
// Parameters:
// - mcpResult: A pointer to an mcp.CallToolResult object containing the result data.
//
// Returns:
// - string: The extracted text content.
// - error: An error if the extraction fails, which can occur in the following cases:
//   - The content slice is empty.
//   - The first element of the content slice is not a map.
//   - The map does not contain a "text" key.
//   - The value associated with the "text" key is not a string.
func getTextFromResult(mcpResult *mcp.CallToolResult) (string, error) {
	if len(mcpResult.Content) == 0 {
		return "", &SSEResultExtractionError{Message: "content is empty"}
	}

	content, ok := mcpResult.Content[0].(mcp.TextContent)
	if !ok {
		return "", &SSEResultExtractionError{Message: "content[0] is not TextContent"}
	}

	return content.Text, nil
}

// Static Resources
func (c *Client) ListResources() (llm.Resources, error) {
	listRequest := mcp.ListResourcesRequest{}
	mcpResources, err := c.mcpClient.ListResources(c.ctx, listRequest)
	if err != nil {
		return nil, &SSEListResourcesError{Message: fmt.Sprintf("Failed to list resources: %v", err)}
	}
	// Convert mcp.Resources to llm.Resources
	resources := llm.Resources{}
	for _, mcpResource := range mcpResources.Resources {
		resources = append(resources, llm.Resource{
			Name:        mcpResource.Name,
			MIMEType:    mcpResource.MIMEType,
			URI:         mcpResource.URI,
			Description: mcpResource.Description,
		})
	}

	return resources, nil
}

type ResourceResult struct {
	//Contents []map[string]interface{}
	Contents []string
}

func (c *Client) ReadResource(uri string) (ResourceResult, error) {
	// Create a request to read the resource
	readRequest := mcp.ReadResourceRequest{}
	readRequest.Params.URI = uri

	// Call the ReadResource API
	mcpResourceResult, err := c.mcpClient.ReadResource(c.ctx, readRequest)
	if err != nil {
		return ResourceResult{}, &SSEReadResourceError{Message: fmt.Sprintf("Failed to read resource: %v", err)}
	}

	//rs, ok := mcpResourceResult.Contents[0].(mcp.TextResourceContents)

	resourceResult := ResourceResult{}
	for _, content := range mcpResourceResult.Contents {

		contentsMap := content.(mcp.TextResourceContents)
		resourceResult.Contents = append(resourceResult.Contents, contentsMap.Text)
	}

	return resourceResult, nil
}

// PromptResult

func (c *Client) ListPrompts() (llm.Prompts, error) {
	promptsRequest := mcp.ListPromptsRequest{}
	mcpPrompts, err := c.mcpClient.ListPrompts(c.ctx, promptsRequest)

	if err != nil {
		return nil, &SSEListPromptsError{Message: fmt.Sprintf("Failed to list prompts: %v", err)}
	}

	prompts := llm.Prompts{}
	for _, mcpPrompt := range mcpPrompts.Prompts {
		prompt := llm.Prompt{}
		prompt.Name = mcpPrompt.Name
		prompt.Description = mcpPrompt.Description // not used

		for _, argument := range mcpPrompt.Arguments {
			prompt.Arguments = append(prompt.Arguments, llm.Argument{
				Name:        argument.Name,
				Description: argument.Description,
				Required:    argument.Required,
			})
		}

		//fmt.Println("📣 Prompt:", mcpPrompt.Name)
		//fmt.Println("  - Description:", mcpPrompt.Description)
		//fmt.Println("  - Arguments:", mcpPrompt.Arguments)

		prompts = append(prompts, prompt)
	}

	return prompts, nil
}

func (c *Client) GetAndFillPrompt(promptName string, arguments map[string]string) (llm.Prompt, error) {
	// Create a request to read the resource
	promptRequest := mcp.GetPromptRequest{}
	promptRequest.Params.Name = promptName
	promptRequest.Params.Arguments = arguments

	promptResult, err := c.mcpClient.GetPrompt(c.ctx, promptRequest)
	if err != nil {
		return llm.Prompt{}, &SSEGetPromptError{Message: fmt.Sprintf("Failed to get prompt: %v", err)}
	}

	prompt := llm.Prompt{}
	prompt.Name = promptName
	prompt.Description = promptResult.Description

	for _, message := range promptResult.Messages {

		prompt.Messages = append(prompt.Messages, llm.Message{
			Role:    string(message.Role),
			Content: message.Content.(mcp.TextContent).Text,
		})
	}
	// 	Content: message.Content.(map[string]interface{})["text"].(string),

	return prompt, nil
}

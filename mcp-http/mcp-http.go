package mcphttp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/parakeet-nest/parakeet/llm"
	"github.com/parakeet-nest/parakeet/tools"
)

type Client struct {
	mcpClient *client.Client
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
	var mcpClient *client.Client
	var err error
	
	// For now, create HTTP client without auth headers until we resolve the options issue
	// TODO: Add support for bearer token authorization
	mcpClient, err = client.NewStreamableHttpClient(baseURL + "/mcp")
	
	// Note: bearerToken support will be added once we resolve the StreamableHTTPCOption vs ClientOption issue
	if bearerToken != "" {
		// This will be implemented once the API is clarified
		// For now, just log that a token was provided but not used
		fmt.Printf("Warning: Bearer token provided but not yet supported for HTTP transport\n")
	}

	if err != nil {
		return Client{}, &HTTPClientCreationError{Message: fmt.Sprintf("Failed to create client: %v", err)}
	}

	return Client{
		mcpClient: mcpClient,
		BaseURL:   baseURL,
		ctx:       ctx,
	}, nil
}

func (c *Client) Start() error {
	if err := c.mcpClient.Start(c.ctx); err != nil {
		return &HTTPClientStartError{Message: fmt.Sprintf("Failed to start client: %v", err)}
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
		return nil, &HTTPClientInitializationError{Message: fmt.Sprintf("Failed to initialize client: %v", err)}
	}

	return initResult, nil
}

func (c *Client) ListTools() ([]llm.Tool, error) {
	toolsRequest := mcp.ListToolsRequest{}
	mcpTools, err := c.mcpClient.ListTools(c.ctx, toolsRequest)

	if err != nil {
		return nil, &HTTPGetToolsError{Message: fmt.Sprintf("Failed to list tools: %v", err)}
	}
	// Convert mcp.Tools to llm.Tools
	// TODO: Handle errors during conversion
	ollamaTools := tools.ConvertMCPTools(mcpTools.Tools)
	return ollamaTools, nil
}

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
		return CallToolResult{}, &HTTPToolCallError{Message: fmt.Sprintf("Failed to call tool: %v", err)}
	}

	text, err := getTextFromResult(mcpResult)
	if err != nil {
		return CallToolResult{}, &HTTPToolCallError{Message: fmt.Sprintf("Failed to extract text from result: %v", err)}
	}

	return CallToolResult{Text: text, Type: "text"}, nil
}

func (c *Client) Close() error {
	return c.mcpClient.Close()
}

func getTextFromResult(mcpResult *mcp.CallToolResult) (string, error) {
	if len(mcpResult.Content) == 0 {
		return "", &HTTPResultExtractionError{Message: "content is empty"}
	}

	content, ok := mcpResult.Content[0].(mcp.TextContent)
	if !ok {
		return "", &HTTPResultExtractionError{Message: "content[0] is not TextContent"}
	}

	return content.Text, nil
}

// Static Resources
func (c *Client) ListResources() (llm.Resources, error) {
	listRequest := mcp.ListResourcesRequest{}
	mcpResources, err := c.mcpClient.ListResources(c.ctx, listRequest)
	if err != nil {
		return nil, &HTTPListResourcesError{Message: fmt.Sprintf("Failed to list resources: %v", err)}
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
	Contents []string
}

func (c *Client) ReadResource(uri string) (ResourceResult, error) {
	// Create a request to read the resource
	readRequest := mcp.ReadResourceRequest{}
	readRequest.Params.URI = uri

	// Call the ReadResource API
	mcpResourceResult, err := c.mcpClient.ReadResource(c.ctx, readRequest)
	if err != nil {
		return ResourceResult{}, &HTTPReadResourceError{Message: fmt.Sprintf("Failed to read resource: %v", err)}
	}

	resourceResult := ResourceResult{}
	for _, content := range mcpResourceResult.Contents {

		contentsMap := content.(mcp.TextResourceContents)
		resourceResult.Contents = append(resourceResult.Contents, contentsMap.Text)
	}

	return resourceResult, nil
}

func (c *Client) ListPrompts() (llm.Prompts, error) {
	promptsRequest := mcp.ListPromptsRequest{}
	mcpPrompts, err := c.mcpClient.ListPrompts(c.ctx, promptsRequest)

	if err != nil {
		return nil, &HTTPListPromptsError{Message: fmt.Sprintf("Failed to list prompts: %v", err)}
	}

	prompts := llm.Prompts{}
	for _, mcpPrompt := range mcpPrompts.Prompts {
		prompt := llm.Prompt{}
		prompt.Name = mcpPrompt.Name
		prompt.Description = mcpPrompt.Description

		for _, argument := range mcpPrompt.Arguments {
			prompt.Arguments = append(prompt.Arguments, llm.Argument{
				Name:        argument.Name,
				Description: argument.Description,
				Required:    argument.Required,
			})
		}

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
		return llm.Prompt{}, &HTTPGetPromptError{Message: fmt.Sprintf("Failed to get prompt: %v", err)}
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

	return prompt, nil
}
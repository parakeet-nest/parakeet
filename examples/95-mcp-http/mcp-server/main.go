package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	// Create MCP server
	s := server.NewMCPServer(
		"mcp-http",
		"0.0.0",
	)

	// =================================================
	// TOOLS:
	// =================================================

	rollDices := mcp.NewTool("rool_dices",
		mcp.WithDescription("Roll some dices"),
		mcp.WithNumber("nb_dices",
			mcp.Required(),
			mcp.Description("Number of dices to roll"),
		),
		mcp.WithNumber("nb_sides",
			mcp.Required(),
			mcp.Description("Number of sides of the dices"),
		),
	)

	s.AddTool(rollDices, rollDicesHandler)

	// Start the HTTP server
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9090"
	}

	log.Println("MCP StreamableHTTP server is running on port", httpPort)

	server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
	).Start(":" + httpPort)
}

func rollDicesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	nbDices := request.GetInt("nb_dices", 1)
	sides := request.GetInt("nb_sides", 6)

	log.Printf("🎲 Rolling %d dice(s) with %d sides each...\n", nbDices, sides)

	roll := func(n, x int) int {
		if n <= 0 || x <= 0 {
			return 0
		}

		results := make([]int, n)
		sum := 0

		for i := range n {
			roll := rand.Intn(x) + 1 
			results[i] = roll
			sum += roll
		}

		return sum
	}

	// Simulate rolling dice
	result := roll(nbDices, sides)

	return mcp.NewToolResultText("Result: " + strconv.Itoa(nbDices) + " dices with " + strconv.Itoa(sides) + " sides: " + strconv.Itoa(result)), nil

}

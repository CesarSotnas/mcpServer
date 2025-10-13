package main

import (
	"fmt"
	"io"
	"net/http"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

type Content struct {
	Title       string  `json:"title" jsonschema:"required,description=The title to submit"`
	Description *string `json:"description" jsonschema:"description=The description to submit"`
}
type MyFunctionsArguments struct {
	Submitter string  `json:"submitter" jsonschema:"required,description=The name of the thing calling this tool (openai, google, claude, etc)"`
	Content   Content `json:"content" jsonschema:"required,description=The content of the message"`
}

func main() {
	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())
	err := server.RegisterTool("hello", "Say hello to a person", func(arguments MyFunctionsArguments) (*mcp_golang.ToolResponse, error) {
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Hello, %s!", arguments.Submitter))), nil
	})
	if err != nil {
		panic(err)
	}

	// Dad Jokes como RESOURCE (dados do banco/API)
	err = server.RegisterResource("jokes://random", "dad_jokes", "Random dad jokes from database", "application/json", func() (*mcp_golang.ResourceResponse, error) {
		// Buscar dados do banco/API
		response, err := http.Get("https://jokefather.com/api/jokes/random")
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar piada: %v", err)
		}
		defer response.Body.Close()

		// Ler o corpo da resposta
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler resposta: %v", err)
		}

		// Retornar como recurso embedded
		jokeText := string(body)
		return mcp_golang.NewResourceResponse(mcp_golang.NewTextEmbeddedResource("jokes://random", jokeText, "application/json")), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterPrompt("prompt_test", "This is a test prompt", func(arguments Content) (*mcp_golang.PromptResponse, error) {
		return mcp_golang.NewPromptResponse("description", mcp_golang.NewPromptMessage(mcp_golang.NewTextContent(fmt.Sprintf("Hello, %s!", arguments.Title)), mcp_golang.RoleUser)), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterResource("test://resource", "resource_test", "This is a test resource", "application/json", func() (*mcp_golang.ResourceResponse, error) {
		return mcp_golang.NewResourceResponse(mcp_golang.NewTextEmbeddedResource("test://resource", "This is a test resource", "application/json")), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.RegisterResource("file://app_logs", "app_logs", "The app logs", "text/plain", func() (*mcp_golang.ResourceResponse, error) {
		return mcp_golang.NewResourceResponse(mcp_golang.NewTextEmbeddedResource("file://app_logs", "This is a test resource", "text/plain")), nil
	})
	if err != nil {
		panic(err)
	}

	err = server.Serve()
	if err != nil {
		panic(err)
	}
}

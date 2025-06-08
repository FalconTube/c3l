package utils

import (
	"context"
	"os"

	"github.com/ollama/ollama/api"
)

type OllamaClient struct {
	Client  *api.Client
	Version string
}

func InitOllamaClient(ollamaHost string) (OllamaClient, error) {
	err := os.Setenv("OLLAMA_HOST", ollamaHost)
	if err != nil {
		return OllamaClient{}, err
	}
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return OllamaClient{}, err
	}

	ctx := context.Background()
	version, err := client.Version(ctx)
	if err != nil {
		return OllamaClient{}, err
	}
	oc := OllamaClient{Client: client, Version: version}

	return oc, nil

}

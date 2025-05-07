package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/ollama/ollama/api"
	"github.com/yarlson/pin"
)

type Cli struct {
	Prompt  string `arg:"" help:"foo"`
	Think   bool   `help:"foo" negatable:""`
	Print   bool   `help:"foo" negatable:""`
	Replace bool   `help:"foo" negatable:""`
}

func (c *Cli) Run() error {

	content := readClipboard()

	// Spinner
	p := pin.New("Waiting for answer...",
		pin.WithSpinnerColor(pin.ColorCyan),
		pin.WithTextColor(pin.ColorYellow),
	)
	cancel := p.Start(context.Background())
	defer cancel()

	// Send prompt and clip content to ollama
	prompt := preparePrompt(c.Prompt, content, c.Think)
	response, err := askOllama(prompt, "qwen3:0.6b")
	if err != nil {
		log.Fatal(err)
	}

	response = trimResponse(response)
	p.Stop("Done!")

	if c.Print {
		fmt.Println(response)
	}

	if c.Replace {
		clipboard.WriteAll(response)
	}

	return nil
}

func askOllama(prompt, model string) (string, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return "", err
	}

	req := &api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
		// set streaming to false
		Stream: new(bool),
	}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		ctx = context.WithValue(ctx, "response", resp.Response)
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		return "", err
	}

	response := ctx.Value("response").(string)
	return response, nil
}

func preparePrompt(prompt, content string, noThink bool) string {
	if noThink {
		prompt = fmt.Sprintf(" /no_think %s", prompt)
	}

	prompt = fmt.Sprintf(
		"Prompt: %s\nContent: %s", prompt, content)
	return prompt
}

func trimResponse(response string) string {
	res := strings.Replace(response, "<think>", "", 1)
	res = strings.Replace(res, "</think>", "", 1)
	res = strings.TrimSpace(res)

	return res
}

func readClipboard() string {
	clipContent, _ := clipboard.ReadAll()
	content := strings.TrimSpace(string(clipContent))
	return content
}

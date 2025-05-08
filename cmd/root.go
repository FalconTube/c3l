package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gen2brain/beeep"
	"github.com/ollama/ollama/api"
	"github.com/yarlson/pin"
)

type Cli struct {
	Prompt  string `arg:"" help:"Prompt being sent to Ollama"`
	Think   bool   `short:"t" help:"If true, uses thinking mode, if applicable in model. If false, adds '/no_think' to prompt" negatable:""`
	Print   bool   `short:"p" help:"If true, prints response to stdout (default: true)" negatable:""`
	Replace bool   `short:"r" help:"If true, put Ollama output on clipboard" negatable:""`
	Model   string `short:"m" help:"Ollama model to use. Available models: https://ollama.com/library" default:"qwen3:0.6b"`
	Notify  bool   `short:"n" help:"If true, display tray notification when finished." default:"false"`
}

func (c *Cli) Run() error {

	content := readClipboard()

	// Spinner
	p := pin.New("Running...",
		pin.WithSpinnerColor(pin.ColorCyan),
		pin.WithTextColor(pin.ColorYellow),
		pin.WithPrefix(fmt.Sprintf("🤖%s", c.Model)),
		pin.WithPrefixColor(pin.ColorMagenta),
	)
	cancel := p.Start(context.Background())
	defer cancel()

	if c.Think == false {
		p.UpdateMessage("Running in no-think mode...")
	}

	// Send prompt and clip content to ollama
	prompt := preparePrompt(c.Prompt, content, c.Think)
	response, err := askOllama(prompt, c.Model)
	if err != nil {
		log.Fatal(err)
	}

	response = trimResponse(response)
	p.Stop("Done!")

	if c.Notify {
		beeep.Notify("Clipllama", "Finished!", "./assets/logo.svg")
	}

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
	if noThink == false {
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

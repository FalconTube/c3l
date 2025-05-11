package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	utils "github.com/FalconTube/c3l/utils"
	"github.com/atotto/clipboard"
	"github.com/gen2brain/beeep"
	"github.com/ollama/ollama/api"
)

type DoCmd struct {
	Prompt string `arg:"" help:"Prompt being sent to Ollama"`
	Flags
}

type Flags struct {
	Think      bool   `short:"t" help:"If true, uses thinking mode, if applicable in model. If false, adds '/no_think' to prompt" negatable:""`
	Print      bool   `short:"p" help:"If true, prints response to stdout (default: true)" negatable:""`
	Replace    bool   `short:"r" help:"If true, put Ollama output on clipboard" negatable:""`
	Model      string `short:"m" help:"Ollama model to use. Available models: https://ollama.com/library" default:"qwen3:0.6b"`
	Notify     bool   `short:"n" help:"If true, display tray notification when finished." negatable:"" default:"false"`
	Expand     bool   `short:"e" help:"Expand given prompt into long version, as defined in $HOME/.c3l.toml " negatable:"" `
	OllamaHost string `help:"IP Address for the Ollama server." env:"OLLAMA_HOST" default:"127.0.0.1:11434"`
}

func (c *DoCmd) Run() error {

	if c.Expand {
		c.Prompt = utils.ExpandPromptFromToml(c.Prompt)
		utils.Logger.Infof(`Expanded prompt to: "%s"`, c.Prompt)
	}

	content := readClipboard()

	// Spinner
	p := utils.InitSpinner(c.Model)
	cancel := p.Start(context.Background())
	defer cancel()

	if c.Think == false {
		p.UpdateMessage("Running in no-think mode...")
	}

	// Send prompt and clip content to ollama
	prompt := preparePrompt(c.Prompt, content, c.Think)
	response, err := askOllama(prompt, c.Model, c.OllamaHost)
	if err != nil {
		utils.Logger.Fatal(err)
	}

	response = trimResponse(response)
	p.Stop("Done!")

	postResponseActions(response, c)

	return nil
}

func askOllama(prompt, model, ollamaHost string) (string, error) {
	os.Setenv("OLLAMA_HOST", ollamaHost)
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

func postResponseActions(response string, c *DoCmd) {
	if c.Notify {
		beeep.Notify("Clipllama", "Finished!", "./assets/logo.svg")
	}

	if c.Print {
		fmt.Println(response)
	}

	if c.Replace {
		clipboard.WriteAll(response)
	}
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

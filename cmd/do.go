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
	utils.Flags
}

type response string

func (c *DoCmd) Run() error {

	if c.Expand {
		expandedPrompt, err := utils.ExpandPromptFromToml(c.Prompt)
		if err != nil {
			return err
		}
		c.Prompt = expandedPrompt
		utils.Logger.Infof(`Expanded prompt to: "%s"`, c.Prompt)
	}

	content, err := readClipboard()
	if err != nil {
		return err
	}
	utils.Logger.Info(content)

	// Spinner
	p := utils.InitSpinner(c.Model)
	cancel := p.Start(context.Background())
	defer cancel()

	if !c.Think {
		p.UpdateMessage("Running in no-think mode...")
	}

	// Send prompt and clip content to ollama
	prompt := preparePrompt(c.Prompt, content, c.Think)
	response, err := askOllama(prompt, c.Model, c.OllamaHost)
	if err != nil {
		return err
	}

	response = trimResponse(response)
	p.Stop("Done!")

	postResponseActions(response, c)

	return nil
}

func askOllama(prompt, model, ollamaHost string) (string, error) {
	err := os.Setenv("OLLAMA_HOST", ollamaHost)
	if err != nil {
		return "", err
	}
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
		ctx = context.WithValue(ctx, response("response"), resp.Response)
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		return "", err
	}

	response := ctx.Value(response("response")).(string)
	return response, nil
}

func preparePrompt(prompt, content string, noThink bool) string {
	if !noThink {
		prompt = fmt.Sprintf(" /no_think %s", prompt)
	}

	prompt = fmt.Sprintf(
		"Prompt: %s\nContent: %s", prompt, content)
	return prompt
}

func postResponseActions(response string, c *DoCmd) {
	if c.Notify {
		_ = beeep.Notify("Clipllama", "Finished!", "./assets/logo.svg")
	}

	if !c.Silent {
		fmt.Println(response)
	}

	if c.Replace {
		_ = clipboard.WriteAll(response)
	}
}

func trimResponse(response string) string {
	res := strings.Replace(response, "<think>", "", 1)
	res = strings.Replace(res, "</think>", "", 1)
	res = strings.TrimSpace(res)

	return res
}

func readClipboard() (string, error) {
	clipContent, err := clipboard.ReadAll()
	if err != nil {
		return "", err
	}

	content := strings.TrimSpace(string(clipContent))
	return content, nil
}

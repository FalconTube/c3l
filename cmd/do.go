package cmd

import (
	"context"
	"fmt"
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

	oc, err := utils.InitOllamaClient(c.OllamaHost)
	if err != nil {
		return err
	}

	if c.Expand {
		expandedPrompt, err := utils.ExpandPromptFromToml(c.Prompt)
		if err != nil {
			return err
		}
		if expandedPrompt != c.Prompt {
			c.Prompt = expandedPrompt
			utils.Logger.Infof(`Expanded prompt to: "%s"`, c.Prompt)
		}

		if c.System != "" {
			expandedSystem, err := utils.ExpandSystemFromToml(c.System)
			if err != nil {
				return err
			}
			if expandedPrompt != c.System {
				c.System = expandedSystem
				utils.Logger.Infof(`Expanded system prompt to: "%s"`, c.System)
			}
		}
	}

	content, err := readClipboard()
	if err != nil {
		return err
	}
	utils.Logger.Debugf("Clipboard content: %s", content)

	// Spinner
	p := utils.InitSpinner(c.Model)
	cancel := p.Start(context.Background())
	defer cancel()

	if !c.Think {
		p.UpdateMessage("Running in no-think mode...")
	}

	// Send prompt and clip content to ollama
	prompt := preparePrompt(oc, c.Prompt, content, c.Think)
	response, err := askOllama(oc, prompt, c.Model, c.Think, c.System)
	if err != nil {
		return err
	}

	response = trimResponse(response)
	p.Stop("Done!")

	postResponseActions(response, c)

	return nil
}

func boolPointer(b bool) *bool {
	return &b
}

func askOllama(ollamaClient utils.OllamaClient, prompt string, model string, think bool, system string) (string, error) {

	req := &api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
		System: system,
		// set streaming to false
		Stream: boolPointer(false),
		// As of Ollama v0.9.0, can set Think in API.
		// Still need to pass it via prompt, if older version in use
		Think: boolPointer(think),
	}

	ctx := context.Background()

	respFunc := func(resp api.GenerateResponse) error {
		ctx = context.WithValue(ctx, response("response"), resp.Response)
		return nil
	}

	err := ollamaClient.Client.Generate(ctx, req, respFunc)
	if err != nil {
		return "", err
	}

	response := ctx.Value(response("response")).(string)
	utils.Logger.Debugf("Response: %s", response)
	return response, nil
}

func preparePrompt(ollamaClient utils.OllamaClient, prompt string, content string, think bool) string {
	// TODO: Proper compare lookup
	if ollamaClient.Version != "0.9.0" {
		if !think {
			prompt = fmt.Sprintf(" /no_think %s", prompt)
		}
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

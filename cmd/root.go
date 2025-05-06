package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

var noThink bool
var doPrint bool
var doReplace bool

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
	fmt.Println(prompt)
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

var rootCmd = &cobra.Command{
	Use:     "clipllama <prompt>",
	Aliases: []string{"cl"},
	Short:   "short",
	Long:    `long`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		userPrompt := args[0]
		content := readClipboard()

		prompt := preparePrompt(userPrompt, content, noThink)
		response, err := askOllama(prompt, "qwen3:0.6b")
		if err != nil {
			log.Fatal(err)
		}

		response = trimResponse(response)

		if doPrint {
			fmt.Println(response)
		}

		if doReplace {
			clipboard.WriteAll(response)
		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&noThink, "nothink", "t", true, `If true, use "/no_think" mode`)
	rootCmd.Flags().BoolVarP(&doReplace, "replace", "r", true, `If true, use replace clipboard content with ollama response`)
	rootCmd.Flags().BoolVarP(&doPrint, "print", "p", false, `If true, prints ollama response to STDOUT`)
}

package main

import (
	cmd "github.com/FalconTube/c3l/cmd"
	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
)

var cli cmd.Cli

func main() {
	// Load CLI
	opt := kong.Configuration(kongyaml.Loader, []string{"~/.c3l.yaml"}...)
	ctx := kong.Parse(&cli,
		kong.Name("c3l"),
		kong.Description("Takes the clipboard content + given prompt and sends it to Ollama"),
		kong.UsageOnError(),
		opt,
	)
	_, err := kong.New(&cli, opt)
	ctx.FatalIfErrorf(err)
	// Run main command
	cli.Run()

}

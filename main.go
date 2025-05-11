package main

import (
	cmd "github.com/FalconTube/c3l/cmd"
	"github.com/alecthomas/kong"
	kongtoml "github.com/alecthomas/kong-toml"
)

type Cli struct {
	Do     cmd.DoCmd     `cmd:"" help:"Send <prompt> and clipboard content to Ollama" aliases:"exec,ask,run"`
	Config cmd.ConfigCmd `cmd:"" help:"Interact with default config at $HOME/.c3l.toml"`
}

var cli Cli

func main() {
	// Load CLI
	// opt := kong.Configuration(kongyaml.Loader, []string{"~/.c3l.yaml"}...)
	opt := kong.Configuration(kongtoml.Loader, []string{"~/.c3l.toml"}...)
	ctx := kong.Parse(&cli,
		kong.Name("c3l"),
		kong.Description("Takes the clipboard content + given prompt and sends it to Ollama"),
		kong.UsageOnError(),
		opt,
	)
	_, err := kong.New(&cli, opt)
	ctx.FatalIfErrorf(err)
	// Run main command
	ctx.Run()
	// cli.Run()

}

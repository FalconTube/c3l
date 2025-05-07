package main

import (
	cmd "github.com/FalconTube/clipllama/cmd"
	"github.com/alecthomas/kong"
	kongyaml "github.com/alecthomas/kong-yaml"
)

var cli cmd.Cli

func main() {
	// Load CLI
	opt := kong.Configuration(kongyaml.Loader, []string{"~/.clipllama.yaml"}...)
	ctx := kong.Parse(&cli, kong.UsageOnError(), opt)
	_, err := kong.New(&cli, opt)
	ctx.FatalIfErrorf(err)
	// Run main command
	cli.Run()

}

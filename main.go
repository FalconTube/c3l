package main

import (
	"fmt"
	"os"

	cmd "github.com/FalconTube/c3l/cmd"
	"github.com/alecthomas/kong"
	kongtoml "github.com/alecthomas/kong-toml"
)

type Cli struct {
	Version VersionFlag `help:"Show version"`

	Do     cmd.DoCmd     `cmd:"" help:"Perform action"`
	Config cmd.ConfigCmd `cmd:"" help:"Perform action"`
}
type VersionFlag bool

var cli Cli
var version string

func (v VersionFlag) BeforeApply() error {
	if version == "" {
		version = "unversioned"
	}
	fmt.Print(version)
	os.Exit(0)
	return nil
}

func main() {
	// Load CLI
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
	_ = ctx.Run()

}

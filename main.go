package main

import (
	"fmt"
	"os"
	"strings"

	cmd "github.com/FalconTube/c3l/cmd"
	"github.com/alecthomas/kong"
	kongtoml "github.com/alecthomas/kong-toml"
)

type Cli struct {
	Version VersionFlag `help:"Show version"`

	// Do is executed when no sucommand is given
	Do      cmd.DoCmd     `cmd:"" default:"withargs" help:"Send <prompt> and clipboard content to Ollama" aliases:"exec,ask,run"`
	Config  cmd.ConfigCmd `cmd:"" help:"Interact with default config at $HOME/.c3l.toml"`
	Prompts cmd.PromptCmd `cmd:"" help:"Interact with prompts"`
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
	// If no args given, print main help
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}
	// first := "Takes the clipboard content + given prompt and sends it to Ollama."
	// sec := "If no subcommand is given, executes the 'do' command."
	// desc := fmt.Sprintf("%s\n %s", first, sec)

	desc := `
	Takes the clipboard content + given prompt and sends it to Ollama.

If no subcommand is given, executes the 'do' command.
	
Examples:
	$ c3l "let's talk about the clipboard content" -p

	$ c3l do "let's talk about the clipboard content" -p

	$ c3l config list

	$ c3l prompts add "let" "let's talk about the clipboard content"
		`
	comp := strings.ReplaceAll(desc, "\n\n", " \n")

	// Load CLI
	opt := kong.Configuration(kongtoml.Loader, []string{"~/.c3l.toml"}...)
	ctx := kong.Parse(&cli,
		kong.Name("c3l"),
		kong.Description(comp),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true, NoExpandSubcommands: true, FlagsLast: true}),
		opt,
	)
	_, err := kong.New(&cli, opt)
	ctx.FatalIfErrorf(err)
	// Run main command
	_ = ctx.Run()

}

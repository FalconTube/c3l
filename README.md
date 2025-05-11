<div align="center">
  <img src="./assets/logo.svg" alt="Clipllama logo" width="150">
</div>

# C3L - Clipllama

Make your clipboard interact with an Ollama server.

Simply use **Copy** in any application and automatically attach it to your prompt!

## Features

- Send a prompt to Ollama and automatically attach the clipboard content
- (Optionally) Update clipboard content with Ollama response
- (Optionally) Print to `stdout`
- Enable/Disable **thinking** mode of supported models like **qwen3**

## Installation

```bash
go install github.com/FalconTube/c3l@latest
```

This will install the latest version of `c3l` on your system and put the `c3l` binary on your path.

## Usage

```bash
$ c3l --help
Usage: c3l <prompt> [flags]

Takes the clipboard content + given prompt and sends it to Ollama

Arguments:
  <prompt>    Prompt being sent to Ollama

Flags:
  -h, --help                             Show context-sensitive help.
  -t, --[no-]think                       If true, uses thinking mode, if applicable in model.
                                         If false, adds '/no_think' to prompt
  -p, --[no-]print                       If true, prints response to stdout (default: true)
  -r, --[no-]replace                     If true, put Ollama output on clipboard
  -m, --model="qwen3:0.6b"               Ollama model to use. Available models:
                                         https://ollama.com/library
  -n, --[no-]notify                      If true, display tray notification when finished.
  -e, --[no-]expand                      Expand given prompt into long version, as defined in
                                         $HOME/.c3l.toml
      --ollama-host="127.0.0.1:11434"    IP Address for the Ollama server ($OLLAMA_HOST)
```

package utils

import (
	"os"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func init() {
	Logger = log.New(os.Stderr)
}

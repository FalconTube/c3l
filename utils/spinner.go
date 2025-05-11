package utils

import (
	"fmt"
	"github.com/yarlson/pin"
)

func InitSpinner(model string) *pin.Pin {
	// Spinner
	p := pin.New("Running...",
		pin.WithSpinnerColor(pin.ColorCyan),
		pin.WithTextColor(pin.ColorYellow),
		pin.WithPrefix(fmt.Sprintf("🤖%s", model)),
		pin.WithPrefixColor(pin.ColorMagenta),
	)
	return p
}

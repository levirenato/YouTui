package ui

import (
	"fmt"
	"os/exec"
	"strings"
)

func copyToClipboard(text string) error {
	copiers := []struct {
		name string
		args []string
	}{
		{"wl-copy", nil},
		{"xclip", []string{"-selection", "clipboard"}},
		{"xsel", []string{"--clipboard", "--input"}},
	}

	for _, c := range copiers {
		path, err := exec.LookPath(c.name)
		if err != nil {
			continue
		}
		cmd := exec.Command(path, c.args...)
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %w", c.name, err)
		}
		return nil
	}

	return fmt.Errorf("no clipboard tool found (install wl-copy, xclip, or xsel)")
}

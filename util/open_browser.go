package util

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenInBrowser tries to open the given URL in the default browser.
func OpenInBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		// rundll32 url.dll,FileProtocolHandler <url>
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		// linux/bsd
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("open browser failed: %w", err)
	}
	return nil
}

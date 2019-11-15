package javascript

import (
	"os"
	"os/exec"
)

func install() error {
	cmd := exec.Command("npm", "ci")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

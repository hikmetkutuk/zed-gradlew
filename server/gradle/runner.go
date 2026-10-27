package gradle

import (
	"os/exec"
)

func Run(workspaceRoot, taskName string) error {
	cmd := exec.Command("./gradlew", taskName, "--console=plain")
	cmd.Dir = workspaceRoot
	return cmd.Run()
}

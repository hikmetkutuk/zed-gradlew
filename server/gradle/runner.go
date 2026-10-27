package gradle

import (
	"fmt"
	"os/exec"
)

func Run(workspaceRoot, taskName string) error {
	script := fmt.Sprintf(`
tell application "Zed"
	activate
end tell
delay 0.4
tell application "System Events"
	tell process "Zed"
		key code 50 using {control down}
		delay 0.8
		keystroke "cd '%s' && ./gradlew %s"
		key code 36
	end tell
end tell`, workspaceRoot, taskName)
	return exec.Command("osascript", "-e", script).Run()
}

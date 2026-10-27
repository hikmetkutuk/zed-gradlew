package gradle

import (
	"bufio"
	"os/exec"
	"strings"
)

type Task struct {
	Name        string
	Description string
	Group       string
}

func Discover(workspaceRoot string) ([]Task, error) {
	cmd := exec.Command("./gradlew", "tasks", "--all", "--console=plain")
	cmd.Dir = workspaceRoot

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parse(string(out)), nil
}

func parse(output string) []Task {
	var tasks []Task
	var currentGroup string

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasSuffix(line, " tasks") && !strings.Contains(line, "-") {
			currentGroup = strings.TrimSuffix(line, " tasks")
			continue
		}

		if len(line) == 0 || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "=") {
			continue
		}

		parts := strings.SplitN(line, " - ", 2)
		name := strings.TrimSpace(parts[0])
		if name == "" || strings.Contains(name, " ") {
			continue
		}

		desc := ""
		if len(parts) == 2 {
			desc = strings.TrimSpace(parts[1])
		}

		tasks = append(tasks, Task{
			Name:        name,
			Description: desc,
			Group:       currentGroup,
		})
	}

	return tasks
}

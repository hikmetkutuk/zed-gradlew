package handler

import (
	"fmt"

	"github.com/hikmetkutuk/zed-gradlew/gradle"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

const cmdRun = "gradlew.run"

type Handler struct {
	tasks         []gradle.Task
	workspaceRoot string
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	if params.RootURI != nil {
		h.workspaceRoot = uriToPath(*params.RootURI)
		go h.refreshTasks()
	}

	codelensProvider := true
	return protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			CodeLensProvider: &protocol.CodeLensOptions{
				ResolveProvider: &codelensProvider,
			},
			ExecuteCommandProvider: &protocol.ExecuteCommandOptions{
				Commands: []string{cmdRun},
			},
		},
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    "zed-gradlew",
			Version: ptr("0.1.0"),
		},
	}, nil
}

func (h *Handler) Initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (h *Handler) Shutdown(context *glsp.Context) error {
	return nil
}

func (h *Handler) CodeLens(context *glsp.Context, params *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	if len(h.tasks) == 0 {
		return nil, nil
	}

	var lenses []protocol.CodeLens
	for i, task := range h.tasks {
		title := fmt.Sprintf("▶ %s", task.Name)
		if task.Description != "" {
			title = fmt.Sprintf("▶ %s  (%s)", task.Name, task.Description)
		}
		lenses = append(lenses, protocol.CodeLens{
			Range: protocol.Range{
				Start: protocol.Position{Line: uint32(i), Character: 0},
				End:   protocol.Position{Line: uint32(i), Character: 0},
			},
			Command: &protocol.Command{
				Title:   title,
				Command: cmdRun,
				Arguments: []any{task.Name},
			},
		})
	}
	return lenses, nil
}

func (h *Handler) ExecuteCommand(context *glsp.Context, params *protocol.ExecuteCommandParams) (any, error) {
	if params.Command != cmdRun {
		return nil, nil
	}
	if len(params.Arguments) == 0 {
		return nil, fmt.Errorf("task name required")
	}

	taskName, ok := params.Arguments[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid task name")
	}

	go gradle.Run(h.workspaceRoot, taskName)
	return nil, nil
}

func (h *Handler) refreshTasks() {
	tasks, err := gradle.Discover(h.workspaceRoot)
	if err != nil {
		return
	}
	h.tasks = tasks
}

func uriToPath(uri protocol.URI) string {
	s := string(uri)
	if len(s) > 7 && s[:7] == "file://" {
		return s[7:]
	}
	return s
}

func ptr[T any](v T) *T {
	return &v
}

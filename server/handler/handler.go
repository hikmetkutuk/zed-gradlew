package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hikmetkutuk/zed-gradlew/gradle"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var logger = func() *log.Logger {
	f, err := os.OpenFile("/tmp/zed-gradlew.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return log.New(os.Stderr, "[zed-gradlew] ", log.LstdFlags)
	}
	return log.New(f, "[zed-gradlew] ", log.LstdFlags)
}()

const cmdRun = "gradlew.run"
const methodShowMessage = "window/showMessage"

type Handler struct {
	tasks         []gradle.Task
	workspaceRoot string
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	if caps, err := json.Marshal(params.Capabilities); err == nil {
		logger.Printf("ClientCapabilities: %s", string(caps))
	}
	if params.RootURI != nil {
		h.workspaceRoot = uriToPath(*params.RootURI)
		logger.Printf("Initialize: workspaceRoot=%s", h.workspaceRoot)
		h.refreshTasks()
		logger.Printf("Initialize: loaded %d tasks", len(h.tasks))
		go h.writeZedTasksJSON()
	}

	codelensProvider := false
	syncKind := protocol.TextDocumentSyncKindFull
	return protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: boolPtr(true),
				Change:    &syncKind,
			},
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

func (h *Handler) DidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	logger.Printf("DidOpen: uri=%s", params.TextDocument.URI)
	go func() {
		context.Notify("workspace/codeLens/refresh", nil)
		logger.Printf("codeLens/refresh sent")
	}()
	return nil
}

func (h *Handler) DidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	logger.Printf("DidChange: uri=%s", params.TextDocument.URI)
	return nil
}

func (h *Handler) Shutdown(context *glsp.Context) error {
	return nil
}

var topTasks = []string{"clean", "build", "test", "bootRun", "assemble"}

func filterTopTasks(tasks []gradle.Task) []gradle.Task {
	byName := map[string]gradle.Task{}
	for _, t := range tasks {
		if !strings.Contains(t.Name, ":") {
			byName[t.Name] = t
		}
	}
	var result []gradle.Task
	for _, name := range topTasks {
		if t, ok := byName[name]; ok {
			result = append(result, t)
		}
	}
	if t, ok := byName["spotlessApply"]; ok {
		result = append(result, t)
	}
	return result
}

func (h *Handler) CodeLens(context *glsp.Context, params *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	filtered := filterTopTasks(h.tasks)
	logger.Printf("CodeLens called: uri=%s tasks=%d top=%d", params.TextDocument.URI, len(h.tasks), len(filtered))
	if len(filtered) == 0 {
		return nil, nil
	}

	var lenses []protocol.CodeLens
	for _, task := range filtered {
		title := fmt.Sprintf("▶ %s", task.Name)
		if task.Name == "spotlessApply" {
			title = "▶ fmt"
		}
		lenses = append(lenses, protocol.CodeLens{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: 0},
			},
			Command: &protocol.Command{
				Title:     title,
				Command:   cmdRun,
				Arguments: []any{task.Name},
			},
		})
	}
	return lenses, nil
}

func (h *Handler) ExecuteCommand(ctx *glsp.Context, params *protocol.ExecuteCommandParams) (any, error) {
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

	logger.Printf("ExecuteCommand: running task=%s", taskName)

	msgType := protocol.MessageTypeInfo
	ctx.Notify(methodShowMessage, protocol.ShowMessageParams{
		Type:    msgType,
		Message: fmt.Sprintf("▶ gradlew %s başlatıldı", taskName),
	})

	go func() {
		if err := gradle.Run(h.workspaceRoot, taskName); err != nil {
			logger.Printf("Task %s failed: %v", taskName, err)
			ctx.Notify(methodShowMessage, protocol.ShowMessageParams{
				Type:    protocol.MessageTypeError,
				Message: fmt.Sprintf("gradlew %s başarısız: %v", taskName, err),
			})
		} else {
			logger.Printf("Task %s completed", taskName)
			ctx.Notify(methodShowMessage, protocol.ShowMessageParams{
				Type:    msgType,
				Message: fmt.Sprintf("✓ gradlew %s tamamlandı.", taskName),
			})
		}
	}()
	return nil, nil
}

type zedTask struct {
	Label   string   `json:"label"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Cwd     string   `json:"cwd"`
}

func (h *Handler) writeZedTasksJSON() {
	filtered := filterTopTasks(h.tasks)
	if len(filtered) == 0 {
		return
	}

	var tasks []zedTask
	for _, t := range filtered {
		tasks = append(tasks, zedTask{
			Label:   fmt.Sprintf("Gradle: %s", t.Name),
			Command: "./gradlew",
			Args:    []string{t.Name},
			Cwd:     "$ZED_WORKTREE_ROOT",
		})
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		logger.Printf("writeZedTasksJSON marshal error: %v", err)
		return
	}

	zedDir := filepath.Join(h.workspaceRoot, ".zed")
	if err := os.MkdirAll(zedDir, 0755); err != nil {
		logger.Printf("writeZedTasksJSON mkdir error: %v", err)
		return
	}

	dest := filepath.Join(zedDir, "tasks.json")
	if _, err := os.Stat(dest); err == nil {
		logger.Printf("writeZedTasksJSON: %s already exists, skipping", dest)
		return
	}
	if err := os.WriteFile(dest, data, 0644); err != nil {
		logger.Printf("writeZedTasksJSON write error: %v", err)
		return
	}
	logger.Printf("writeZedTasksJSON: wrote %d tasks to %s", len(tasks), dest)
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

func boolPtr(v bool) *bool {
	return &v
}

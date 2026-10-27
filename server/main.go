package main

import (
	"github.com/hikmetkutuk/zed-gradlew/handler"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "zed-gradlew"

func main() {
	h := handler.NewHandler()

	s := server.NewServer(&protocol.Handler{
		Initialize:                   h.Initialize,
		Initialized:                  h.Initialized,
		Shutdown:                     h.Shutdown,
		TextDocumentDidOpen:          h.DidOpen,
		TextDocumentDidChange:        h.DidChange,
		TextDocumentCodeLens:         h.CodeLens,
		WorkspaceExecuteCommand:      h.ExecuteCommand,
		SetTrace: func(context *glsp.Context, params *protocol.SetTraceParams) error {
			return nil
		},
	}, lsName, false)

	s.RunStdio()
}

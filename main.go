package main

import (
	"bufio"
	"dbwf-ls/analysis"
	"dbwf-ls/error"
	"dbwf-ls/jsonrpc"
	"dbwf-ls/lsp"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("[dbwf-ls] HOME NOT SET???")
		log.Fatal(err)
		panic(err)
	}
	logger := getLogger(home + "/.config/dbwf-ls/log.txt")
	logger.Println("Here comes the crescendo!!")
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(jsonrpc.Split)

	state := analysis.NewState()
	writer := os.Stdout

	for sc.Scan() {
		msg := sc.Bytes()
		method, contents, err := jsonrpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Some errors occurred: %s", err)
			continue
		}
		handleMessage(logger, writer, state, method, contents)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply, _ := jsonrpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
	logger.Printf("Received msg with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitialiseRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("WHY IS IT AN ACORN? %s", err)
		}
		logger.Printf("Attached to %s client version %s", request.Params.ClientInfo.Name, request.Params.ClientInfo.Version)
		msg := lsp.NewInitialiseResponse(request.ID)
		writeResponse(writer, msg)
		logger.Print("Reply sent")
	case "textDocument/didOpen":
		var noti lsp.DidOpenTextNotification
		if err := json.Unmarshal(contents, &noti); err != nil {
			logger.Printf("textDocument/didOpen %s", err)
			return
		}
		notification := state.OpenDocument(noti.Params.TextDocument.URI, noti.Params.TextDocument.Text, logger)
		logger.Printf("Editing %s", noti.Params.TextDocument.URI)
		writeResponse(writer, notification)
		logger.Print("Diagnostics sent")
	case "textDocument/didChange":
		var noti lsp.DidChangeTextNotification
		if err := json.Unmarshal(contents, &noti); err != nil {
			logger.Printf("textDocument/didChange %s", err)
			return
		}

		notification := lsp.PublishDiagnosticsNotification{}

		for _, change := range noti.Params.ContentChanges {
			notification = state.UpdateDocument(noti.Params.TextDocument.URI, change.Text, logger)
		}
		logger.Printf("Updated %s", noti.Params.TextDocument.URI)
		writeResponse(writer, notification)

		logger.Print("Diagnostics sent")
	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/hover %s", err)
			return
		}

		// Hover response
		response, err := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position, logger)
		if err != nil {
			writeResponse(writer, lsp.ErrorResponse{
				Response: lsp.Response{
					RPC: "2.0",
					ID:  &request.ID,
				},
				Error: error.ResponseError{
					Code:    error.InternalError,
					Message: fmt.Sprintf("Internal error: %s\n", err),
				},
			})
		} else {
			writeResponse(writer, response)
		}

		logger.Print("Hover response sent")
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/definition %s", err)
			return
		}

		// Definition response
		response, err := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position, logger)
		if err != nil {
			logger.Println(err)
		} else {
			writeResponse(writer, response)
		}
		logger.Print("Definition response sent")
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/codeAction %s", err)
			return
		}

		// CodeAction response
		response, err := state.CodeAction(request.ID, request.Params.TextDocument.URI, logger)
		if err != nil {
			writeResponse(writer, lsp.ErrorResponse{
				Response: lsp.Response{
					RPC: "2.0",
					ID:  &request.ID,
				},
				Error: error.ResponseError{
					Code:    error.InternalError,
					Message: fmt.Sprintf("Internal error: %s\n", err),
				},
			})
		} else {
			writeResponse(writer, response)
		}
		logger.Print("Code Action response sent")
	case "textDocument/formatting":
		var request lsp.DocumentFormattingRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/formatting %s", err)
			return
		}

		// Formatting response
		response, err := state.DocumentFormatting(request.ID, request.Params.TextDocument.URI, request.Params.Options, logger)

		if err != nil {
			writeResponse(writer, lsp.ErrorResponse{
				Response: lsp.Response{
					RPC: "2.0",
					ID:  &request.ID,
				},
				Error: error.ResponseError{
					Code:    error.InternalError,
					Message: fmt.Sprintf("Internal error: %s", err),
				},
			})
		} else {
			writeResponse(writer, response)
		}
		logger.Print("Formatting response sent")
	case "textDocument/completion":
		var request lsp.CompletionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("textDocument/completion %s", err)
			return
		}

		// Completion response
		response, err := state.Completion(request.ID, request.Params.TextDocument.URI, request.Params.Position, logger)

		if err != nil {
			writeResponse(writer, lsp.ErrorResponse{
				Response: lsp.Response{
					RPC: "2.0",
					ID:  &request.ID,
				},
				Error: error.ResponseError{
					Code:    error.InternalError,
					Message: fmt.Sprintf("Internal error: %s", err),
				},
			})
		} else {
			writeResponse(writer, response)
		}
		logger.Print("Completion response sent")

	}
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic("WHY IS IT AN ACORN?")
	}

	return log.New(logfile, "[dbwf-ls]", log.Ldate|log.Ltime|log.Lshortfile)
}

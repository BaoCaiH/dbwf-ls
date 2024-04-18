package analysis

import (
	"dbwf-ls/lsp"
	"log"
	"regexp"
	"strings"
)

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func (s *State) OpenDocument(uri, text string) lsp.PublishDiagnosticsNotification {
	s.Documents[uri] = text

	diagnostics := diagnose(s.Documents[uri])

	return lsp.PublishDiagnosticsNotification{
		Notification: lsp.Notification{
			RPC:    "2.0",
			Method: "textDocument/publishDiagnostics",
		},
		Params: lsp.PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diagnostics,
		},
	}
}

func (s *State) UpdateDocument(uri, text string) lsp.PublishDiagnosticsNotification {
	s.Documents[uri] = text

	diagnostics := diagnose(s.Documents[uri])

	return lsp.PublishDiagnosticsNotification{
		Notification: lsp.Notification{
			RPC:    "2.0",
			Method: "textDocument/publishDiagnostics",
		},
		Params: lsp.PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diagnostics,
		},
	}
}

func (s *State) Hover(id int, uri string, position lsp.Position, logger *log.Logger) (lsp.HoverResponse, error) {
	document := s.Documents[uri]

	line := strings.Split(document, "\n")[position.Line]

	word, err := wordAtCursor(line, position, logger)
	if err != nil {
		return lsp.HoverResponse{}, err
	}

	// Hover response
	response := lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: word,
		},
	}

	return response, nil
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	// document := s.Documents[uri]

	// Definition response
	response := lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{Line: position.Line - 1, Character: 0},
				End:   lsp.Position{Line: position.Line - 1, Character: 0},
			},
		},
	}

	return response
}

func (s *State) CodeAction(id int, uri string, logger *log.Logger) (lsp.CodeActionResponse, error) {
	document := s.Documents[uri]

	actions := []lsp.CodeAction{}
	re, err := regexp.Compile("\\s+$")
	if err != nil {
		logger.Printf("CodeAction Regexp Compile %s", err)
		return lsp.CodeActionResponse{}, err
	}
	for row, line := range strings.Split(document, "\n") {
		loc := re.FindStringIndex(line)
		if loc != nil {
			dropTrailingWhitespacesEdit := map[string][]lsp.TextEdit{}
			dropTrailingWhitespacesEdit[uri] = []lsp.TextEdit{
				{
					Range:   lsp.LineRange(row, loc[0], loc[1]),
					NewText: "",
				},
			}
			actions = append(actions, lsp.CodeAction{
				Title: "Remove trailing whitespaces",
				Edit:  &lsp.WorkspaceEdit{Changes: dropTrailingWhitespacesEdit},
			})
		}
	}

	// Code Action response
	response := lsp.CodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: actions,
	}

	return response, nil
}

func (s *State) DocumentFormatting(id int, uri string, opts lsp.FormattingOptions, logger *log.Logger) (lsp.DocumentFormattingResponse, error) {
	document := s.Documents[uri]

	// Note: No need to do tabs for yaml, apparently. But I wrote it so I'm keeping it
	tabs := strings.Repeat(" ", opts.TabSize)
	trimTrailingWhitespace, trimFinalNewlines := true, true
	insertSpace, trimTrailingWhitespace, trimFinalNewlines := true, true, true
	if opts.InsertSpaces != nil {
		insertSpace = *opts.InsertSpaces
	}
	if opts.TrimTrailingWhitespace != nil {
		trimTrailingWhitespace = *opts.TrimTrailingWhitespace
	}
	if opts.TrimFinalNewlines != nil {
		trimFinalNewlines = *opts.TrimFinalNewlines
	}

	edits := []lsp.TextEdit{}
	re, err := regexp.Compile("\n+$")
	if err != nil {
		logger.Printf("Formatting Regexp Compile %s", err)
		return lsp.DocumentFormattingResponse{}, err
	}
	loc := re.FindStringIndex(document)
	if loc != nil && trimFinalNewlines {
		trim := loc[1] - loc[0] + 1
		lines := strings.Split(document, "\n")
		rows := len(lines)
		trimFinalNewlinesEdit := lsp.TextEdit{
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      rows - trim,
					Character: len(lines[rows-trim]),
				},
				End: lsp.Position{
					Line:      rows - 1,
					Character: 0,
				},
			},
			NewText: "",
		}
		edits = append(edits, trimFinalNewlinesEdit)
	}

	re, err = regexp.Compile("\\s+$")
	if err != nil {
		logger.Printf("Formatting Regexp Compile %s", err)
		return lsp.DocumentFormattingResponse{}, err
	}
	reTab, err := regexp.Compile("\\t")
	if err != nil {
		logger.Printf("Formatting Regexp Compile %s", err)
		return lsp.DocumentFormattingResponse{}, err
	}
	for row, line := range strings.Split(document, "\n") {
		loc := re.FindStringIndex(line)
		if loc != nil && trimTrailingWhitespace {
			edits = append(edits, lsp.TextEdit{
				Range:   lsp.LineRange(row, loc[0], loc[1]),
				NewText: "",
			})
		}
		locs := reTab.FindAllStringIndex(line, -1)
		if locs != nil && insertSpace {
			logger.Print("Found some open tabs")
			for _, loc := range locs {
				edits = append(edits, lsp.TextEdit{
					Range:   lsp.LineRange(row, loc[0], loc[1]),
					NewText: tabs,
				})
			}
		}
	}

	// Formatting response
	response := lsp.DocumentFormattingResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: edits,
	}

	return response, nil
}

func (s *State) Completion(id int, uri string, position lsp.Position, logger *log.Logger) (lsp.CompletionResponse, error) {
	document := s.Documents[uri]

	items := []lsp.CompletionItem{}
	line := strings.Split(document, "\n")[position.Line]
	word, err := wordAtCursor(line, position, logger)
	if err != nil {
		return lsp.CompletionResponse{}, err
	}

	leading, err := leadingSpaces(line, logger)
	if err != nil {
		return lsp.CompletionResponse{}, err
	}

	items = append(items, complete(word, leading)...)

	// Completion response
	response := lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.CompletionList{
			IsIncomplete: true,
			Items:        items,
		},
	}

	return response, nil
}

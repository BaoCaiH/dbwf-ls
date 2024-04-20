package analysis

import (
	"dbwf-ls/lsp"
	"errors"
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

// Handler for when document opened
// It simply add the full document to the current state
// and provide diagnostics
// For task diagnostics, it needs a `description` right below the `task_key` to be discovered
func (s *State) OpenDocument(uri, text string, logger *log.Logger) lsp.PublishDiagnosticsNotification {
	s.Documents[uri] = text

	diagnostics := diagnose(s.Documents[uri], logger)

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

// Handler for when document changed
// It also add the full document to the current state
// and provide diagnostics
// For task diagnostics, it needs a `description` right below the `task_key` to be discovered
func (s *State) UpdateDocument(uri, text string, logger *log.Logger) lsp.PublishDiagnosticsNotification {
	s.Documents[uri] = text

	diagnostics := diagnose(s.Documents[uri], logger)

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

// Handler for hover request
// Selected keywords in `Keywords` are filled with documentations from databricks
func (s *State) Hover(id int, uri string, position lsp.Position, logger *log.Logger) (lsp.HoverResponse, error) {
	document := s.Documents[uri]

	line := strings.Split(document, "\n")[position.Line]

	word, err := wordAtCursor(line, position, logger)
	if err != nil {
		return lsp.HoverResponse{}, err
	}

	content := lsp.MarkupContent{}
	if Keywords[word].hover == content {
		content = lsp.MarkupContent{
			Kind:  "plaintext",
			Value: "No information available",
		}
	} else {
		content = Keywords[word].hover
	}

	// Hover response
	response := lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: content,
		},
	}

	return response, nil
}

// Handler for go to definition request
// Parse the document and find where a task or a cluster is defined
// For task, it needs a `description` right below the `task_key` to be discovered
func (s *State) Definition(id int, uri string, position lsp.Position, logger *log.Logger) (lsp.DefinitionResponse, error) {
	document := s.Documents[uri]
	lines := strings.Split(document, "\n")
	item_name, err := wordAtCursor(lines[position.Line], position, logger)
	if err != nil {
		return lsp.DefinitionResponse{}, err
	}
	if item_name == "" {
		return lsp.DefinitionResponse{}, errors.New("Not a word")
	}
	matched, err := regexp.MatchString("job_cluster_key:|task_key:", lines[position.Line])
	if err != nil {
		return lsp.DefinitionResponse{}, err
	}
	if !matched {
		return lsp.DefinitionResponse{}, errors.New("Not task or cluster")
	}

	item := findDefinition(document, item_name, logger)

	if item.defined == lsp.LineRange(0, 0, 0) {
		return lsp.DefinitionResponse{}, errors.New("Not defined")
	}

	// Definition response
	response := lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			URI:   uri,
			Range: item.defined,
		},
	}

	return response, nil
}

// Handler for code action request
// For now it does the same thing as the simplest format
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

// Handler for format request
// It can insert spaces, trim whitespaces and trailing new lines
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

// Handler for completion request
// Selected keywords in `Keywords` are filled with examples
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

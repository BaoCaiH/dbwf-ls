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

func diagnose(document string) []lsp.Diagnostics {
	diagnostics := []lsp.Diagnostics{}
	foundHealth := false
	for _, line := range strings.Split(document, "\n") {
		matched, err := regexp.MatchString("^health:\\s*?/??/??.*?$", line)
		if err != nil {
			return diagnostics
		}
		if matched {
			foundHealth = true
		}
	}

	if !foundHealth {
		diagnostics = append(diagnostics, lsp.Diagnostics{
			Range:    lsp.LineRange(0, 0, 0),
			Severity: 2,
			Source:   "dbwf-ls",
			Message:  "Please consider adding `health` check to workflows.",
		})
	}

	return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostics {
	s.Documents[uri] = text

	return diagnose(s.Documents[uri])
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostics {
	s.Documents[uri] = text

	return diagnose(s.Documents[uri])
}

func wordAtCursor(line string, position lsp.Position, logger *log.Logger) (string, error) {
	re, err := regexp.Compile("\\W")
	if err != nil {
		logger.Printf("Regexp Compile %s", err)
		return "", err
	}

	// Because the flocking cursor is 1 step ahead of the line while typing
	// So this can fail, quietly, damn.
	char := position.Character
	if char == len(line) {
		char--
	}

	if loc := re.FindStringIndex(line[char : char+1]); loc != nil {
		return "", nil
	}

	start, end := 0, 0
	if locs := re.FindAllStringIndex(line, -1); locs != nil {
		for _, loc := range locs {
			if loc[0] > position.Character {
				end = loc[0]
			}
			if loc[1] <= position.Character {
				start = loc[1]
			}
			if end != 0 {
				break
			}
		}
	}
	if end == 0 {
		end = len(line)
	}

	return line[start:end], nil
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
	// tabs := strings.Repeat(" ", opts.TabSize*2)
	trimTrailingWhitespace, trimFinalNewlines := true, true
	// insertSpace, trimTrailingWhitespace, trimFinalNewlines := true, true, true
	// if opts.InsertSpaces != nil {
	// 	insertSpace = *opts.InsertSpaces
	// }
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
	// reTab, err := regexp.Compile("\\t")
	// if err != nil {
	// 	logger.Printf("Formatting Regexp Compile %s", err)
	// 	return lsp.DocumentFormattingResponse{}, err
	// }
	for row, line := range strings.Split(document, "\n") {
		loc := re.FindStringIndex(line)
		if loc != nil && trimTrailingWhitespace {
			edits = append(edits, lsp.TextEdit{
				Range:   lsp.LineRange(row, loc[0], loc[1]),
				NewText: "",
			})
		}
		// locs := reTab.FindAllStringIndex(line, -1)
		// if locs != nil && insertSpace {
		// 	logger.Print("Found some open tabs")
		// 	for _, loc := range locs {
		// 		edits = append(edits, lsp.TextEdit{
		// 			Range:   lsp.LineRange(row, loc[0], loc[1]),
		// 			NewText: tabs,
		// 		})
		// 	}
		// }
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
	logger.Print(line, position)
	word, err := wordAtCursor(line, position, logger)
	logger.Print(word)
	if err != nil {
		return lsp.CompletionResponse{}, err
	}
	logger.Print(word)
	items = append(items, lsp.CompletionItem{
		Label:         word,
		Detail:        "Current typing word",
		Documentation: "Nothing to document here",
	})

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

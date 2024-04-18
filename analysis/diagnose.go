package analysis

import (
	"dbwf-ls/lsp"
	"regexp"
	"strings"
)

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

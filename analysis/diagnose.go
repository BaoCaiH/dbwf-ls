package analysis

import (
	"cmp"
	"dbwf-ls/lsp"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"
)

type requiredKeys struct {
	name string
	Diag
	found bool
}

func diagnose(document string, logger *log.Logger) []lsp.Diagnostics {
	diagnostics := []lsp.Diagnostics{}
	required_keywords := map[string]requiredKeys{}
	for k, v := range Keywords {
		if v.diag.severity != 0 {
			required_keywords[k] = requiredKeys{
				Diag:  v.diag,
				found: false,
			}
		}
	}
	logger.Println(required_keywords)
	for _, line := range strings.Split(document, "\n") {
		for k, v := range required_keywords {
			if v.found {
				continue
			}
			matched, err := regexp.MatchString(fmt.Sprintf("^%s:\\s*?/??/??.*?$", k), line)
			if err != nil {
				return diagnostics
			}
			if matched {
				v.found = true
				required_keywords[k] = v
			}
		}
	}

	documentLength := len(strings.Split(document, "\n"))
	for _, v := range required_keywords {
		if !v.found {
			diagnostics = append(diagnostics, lsp.Diagnostics{
				Range:    lsp.LineRange(documentLength-1, 0, 0),
				Severity: v.severity,
				Source:   "dbwf-ls",
				Message:  v.help,
			})
		}
	}

	slices.SortFunc(diagnostics, func(a, b lsp.Diagnostics) int {
		if cmp.Compare(a.Severity, b.Severity) == 0 {
			return cmp.Compare(a.Message, b.Message)
		}
		return cmp.Compare(a.Severity, b.Severity)
	})

	return diagnostics
}

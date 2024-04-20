package analysis

import (
	"dbwf-ls/lsp"
	"strings"
)

func hammingRatio(input, keyword string) float32 {
	compareLength := len(input)
	if len(input) > len(keyword) {
		compareLength = len(keyword)
	}

	dist := 0
	for i := range compareLength {
		if input[i] != keyword[i] {
			dist += 1
		}
	}

	return (float32(compareLength) - float32(dist)) / float32(compareLength)
}

func complete(word, leading string) []lsp.CompletionItem {
	options := []lsp.CompletionItem{}
	for kw, completions := range Keywords {
		if hammingRatio(word, kw) >= 0.75 {
			for _, completion := range completions.completions {
				options = append(options, lsp.CompletionItem{
					Label:         kw,
					Kind:          completion.kind,
					Detail:        completion.detail,
					Documentation: completion.documentation,
					InsertText:    strings.ReplaceAll(completion.insertText, "\n", "\n"+leading),
				})
			}
		}
	}
	return options
}

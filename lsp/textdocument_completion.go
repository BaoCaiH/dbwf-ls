package lsp

type CompletionRequest struct {
	Request
	Params CompletionParams `json:"params"`
}

type CompletionParams struct {
	TextDocumentPositionParams
}

type CompletionResponse struct {
	Response
	Result CompletionList `json:"result"`
}

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

type CompletionItem struct {
	Label         string        `json:"label"`
	Kind          int           `json:"kind"`
	Detail        string        `json:"detail"`
	Documentation MarkupContent `json:"documentation"`
	InsertText    string        `json:"insertText"`
}

type MarkupContent struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

var CompletionItemKind = map[string]int{
	"Text":          1,
	"Method":        2,
	"Function":      3,
	"Constructor":   4,
	"Field":         5,
	"Variable":      6,
	"Class":         7,
	"Interface":     8,
	"Module":        9,
	"Property":      10,
	"Unit":          11,
	"Value":         12,
	"Enum":          13,
	"Keyword":       14,
	"Snippet":       15,
	"Color":         16,
	"File":          17,
	"Reference":     18,
	"Folder":        19,
	"EnumMember":    20,
	"Constant":      21,
	"Struct":        22,
	"Event":         23,
	"Operator":      24,
	"TypeParameter": 25,
}

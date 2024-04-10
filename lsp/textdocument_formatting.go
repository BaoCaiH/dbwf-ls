package lsp

type DocumentFormattingRequest struct {
	Request
	Params DocumentFormattingParams `json:"params"`
}

type DocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

type FormattingOptions struct {
	TabSize                int   `json:"tabSize"`
	InsertSpaces           *bool `json:"insertSpaces,omitempty"`
	TrimTrailingWhitespace *bool `json:"trimTrailingWhitespace,omitempty"`
	InsertFinalNewline     *bool `json:"insertFinalNewline,omitempty"`
	TrimFinalNewlines      *bool `json:"trimFinalNewlines,omitempty"`
}

type DocumentFormattingResponse struct {
	Response
	Result []TextEdit `json:"result"`
}

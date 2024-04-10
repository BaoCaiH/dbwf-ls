package lsp

type InitialiseRequest struct {
	Request
	Params InitialiseRequestParams `json:"params"`
}

type InitialiseRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitialiseResponse struct {
	Response
	Result InitialiseResult `json:"result"`
}

type InitialiseResult struct {
	Capabitities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync           int  `json:"textDocumentSync"`
	HoverProvider              bool `json:"hoverProvider"`
	DefinitionProvider         bool `json:"definitionProvider"`
	CodeActionProvider         bool `json:"codeActionProvider"`
	CompletionProvider         bool `json:"completionProvider"`
	DocumentFormattingProvider bool `json:"documentFormattingProvider"`
}
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitialiseResponse(id int) InitialiseResponse {
	return InitialiseResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitialiseResult{
			Capabitities: ServerCapabilities{
				// TODO: Change to incremental
				TextDocumentSync:           1,
				HoverProvider:              true,
				DefinitionProvider:         true,
				CodeActionProvider:         true,
				DocumentFormattingProvider: true,
			},
			ServerInfo: ServerInfo{
				Name:    "dbwf-ls",
				Version: "timonthy",
			},
		},
	}
}

package error

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

const (
	ParseError           int = -32700
	InvalidRequest           = -32600
	MethodNotFound           = -32601
	InvalidParams            = -32602
	InternalError            = -32603
	ServerNotInitialized     = -32002
	UnknownErrorCode         = -32001
	RequestFailed            = -32803
	ServerCancelled          = -32802
	ContentModified          = -32801
	RequestCancelled         = -32800
)

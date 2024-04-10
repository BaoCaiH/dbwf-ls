package jsonrpc_test

import (
	"dbwf-ls/jsonrpc"
	"testing"
)

type EncodingTest struct {
	Dummy bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 14\r\n\r\n{\"Dummy\":true}"
	actual, err := jsonrpc.EncodeMessage(EncodingTest{Dummy: true})

	if err != nil {
		t.Fatal(err)
	}

	if expected != actual {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	incomingMessage := "Content-Length: 15\r\n\r\n{\"Method\":\"no\"}"
	method, content, err := jsonrpc.DecodeMessage([]byte(incomingMessage))
	if err != nil {
		t.Fatal(err)
	}
	if len(content) != 15 {
		t.Fatalf("Content-Length, expected: 15, got: %d", len(content))
	}
	if string(content) != "{\"Method\":\"no\"}" {
		t.Fatalf("Content, expected: {\"Method\":\"no\"}, got: %s", content)
	}
	if method != "no" {
		t.Fatalf("Method, expected: \"no\", got: %s", method)
	}
}

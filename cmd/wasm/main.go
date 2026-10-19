//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/ZachariahAyers/algokit/v2/internal/demo"
)

type response struct {
	Result *demo.Result `json:"result,omitempty"`
	Error  string       `json:"error,omitempty"`
}

func main() {
	runner := js.FuncOf(run)
	js.Global().Set("algokitRun", runner)
	if ready := js.Global().Get("algokitReady"); ready.Type() == js.TypeFunction {
		ready.Invoke()
	}
	select {}
}

func run(_ js.Value, args []js.Value) any {
	if len(args) != 1 || args[0].Type() != js.TypeString {
		return encode(response{Error: "the demo expects one JSON request"})
	}
	request, err := demo.ParseRequest([]byte(args[0].String()))
	if err != nil {
		return encode(response{Error: err.Error()})
	}
	result := demo.Execute(request)
	return encode(response{Result: &result})
}

func encode(value response) string {
	data, err := json.Marshal(value)
	if err != nil {
		return `{"error":"the result could not be encoded"}`
	}
	return string(data)
}

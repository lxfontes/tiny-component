package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"go.wasmcloud.dev/component"
	"go.wasmcloud.dev/component/net/wasihttp"
)

func main() {
	fmt.Println("Running it!")
}

var logger = component.DefaultLogger

func init() {
	component.Init()

	wasihttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Handling request", slog.String("context", "Handle"))

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("yes"))
	})
}

//go:generate wit-bindgen-go generate --world default --out gen ./wit

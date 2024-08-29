package main

import (
	"log/slog"
	"net/http"

	"go.wasmcloud.dev/component"
	"go.wasmcloud.dev/component/net/wasihttp"
)

func main() {}

var logger = component.DefaultLogger

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Info("Handling request", slog.String("context", "Handle"))

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("yes"))
}

func init() {
	wasihttp.Handle(serveHTTP)
}

//go:generate wit-bindgen-go generate --world default --out gen ./wit

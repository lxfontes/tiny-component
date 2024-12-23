package main

//go:generate go run go.bytecodealliance.org/cmd/wit-bindgen-go generate --world default --out gen ./wit

import (
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"

	"go.wasmcloud.dev/component/log/wasilog"
	"go.wasmcloud.dev/component/net/wasihttp"
)

var (
	logger     = wasilog.DefaultLogger
	httpClient = wasihttp.DefaultClient
)

func init() {
	wasihttp.Handle(http.HandlerFunc(proxyHandler))
}

func blastHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	l := logger.With("context", "Handle")
	l.Info("Handling request")

	rng := rand.IntN(100)
	outgoingLocation := "east"
	if rng < 50 {
		outgoingLocation = "west"
	}

	// WASI Roundtripper
	l.Info("Creating request")
	req, err := http.NewRequest(http.MethodGet, "http://www.randomnumberapi.com/api/v1.0/random?min=100&max=1000&count=5", nil)
	if err != nil {
		http.Error(w, "failed to create request", http.StatusBadGateway)
		return
	}

	l.Info("Executing request")
	resp, err := httpClient.Do(req)
	if err != nil {
		http.Error(w, "failed to make outbound request", http.StatusBadGateway)
		return
	}

	l.Info("Defering body close")
	defer resp.Body.Close()

	l.Info("Checking response status")
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "invalid http status from upstream", http.StatusBadGateway)
		return
	}

	l.Info("Proxying")
	w.Header().Set("X-Outgoing-Location", outgoingLocation)
	w.WriteHeader(http.StatusOK)

	l.Info("Forwarding response")
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		l.Error("Failed to forward response", slog.Any("error", err))
	}
}

func main() {}

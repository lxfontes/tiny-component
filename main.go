package main

//go:generate wit-bindgen-go generate --world default --out gen ./wit

import (
	"fmt"
	"io"
	"log/slog"
	"math/rand/v2"
	"net/http"

	monotonicclock "github.com/lxfontes/tiny-component/gen/wasi/clocks/monotonic-clock"
	"github.com/ydnar/wasm-tools-go/cm"
	"go.wasmcloud.dev/component"
	"go.wasmcloud.dev/component/lattice"
	"go.wasmcloud.dev/component/net/wasihttp"
	"go.wasmcloud.dev/component/time/wasitime"
)

var (
	logger     = component.DefaultLogger
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

	l.Info("Setting link name", slog.String("location", outgoingLocation))
	interfaces := []lattice.CallTargetInterface{
		lattice.NewCallTargetInterface("wasi", "http", "outgoing-handler"),
	}
	lattice.SetLinkName(outgoingLocation, cm.ToList(interfaces))

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
	w.Header().Set("X-Now", wasitime.Now().String())
	w.Header().Set("X-Mono", fmt.Sprintf("%d", monotonicclock.Now()))
	w.WriteHeader(http.StatusOK)

	l.Info("Forwarding response")
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		l.Error("Failed to forward response", slog.Any("error", err))
	}
}

func main() {}

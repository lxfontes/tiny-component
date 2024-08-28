package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/lxfontes/component"
	"github.com/lxfontes/component/net/wasihttp"
)

func main() {
	fmt.Println("Running it!")
}

func init() {
	component.Init()

	wasihttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		client := wasihttp.NewClient()
		req, err := http.NewRequest(http.MethodGet, "http://www.randomnumberapi.com/api/v1.0/random?min=100&max=1000&count=5", nil)
		if err != nil {
			fmt.Println("failed to create request", err)
			os.Exit(1)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("failed to make outbound request", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("expected status code: %d, got: %d\n", http.StatusOK, resp.StatusCode)
			os.Exit(1)
		}

		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("failed to read response body", err)
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(raw)
	})
}

//go:generate wit-bindgen-go generate --world default --out gen ./wit

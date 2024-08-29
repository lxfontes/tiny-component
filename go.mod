module github.com/lxfontes/tiny-component

go 1.22.5

require (
	github.com/ydnar/wasm-tools-go v0.1.5
	go.wasmcloud.dev/component v0.0.0-20240828212418-086462b76175
)

require (
	github.com/samber/lo v1.44.0 // indirect
	github.com/samber/slog-common v0.17.1 // indirect
	golang.org/x/text v0.16.0 // indirect
)

replace go.wasmcloud.dev/component => ../../wasmcloud/component-sdk-go

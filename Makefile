all:
	tinygo build -target=wasip2 --wit-package ./wit --wit-world default -o main.wasm main.go

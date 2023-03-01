.PHONY: wasm
wasm:
	GOOS=js GOARCH=wasm go build -o ./web/yy.wasm

.PHONY: wasmo
wasmo: wasm
	wasm-opt -O ./web/yy.wasm -o ./web/yy.wasm --enable-bulk-memory
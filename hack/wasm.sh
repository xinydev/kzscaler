cd cmd/proxy
tinygo build -o proxy.wasm -scheduler=none -target=wasi main.go

envoy -c ./envoy.yaml --concurrency 2 --log-format '%v'
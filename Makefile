goimports := golang.org/x/tools/cmd/goimports@v0.1.5
golangci_lint := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.42.0

.PHONY: build.wasm
build.wasm:
	#echo ./cmd/proxy/main.go | xargs -Ip tinygo build -o p.wasm -scheduler=none -target=wasi p
	echo ./cmd/proxy/main.go | xargs -Ip tinygo build -o p.wasm -target=wasi p

.PHONY: run
run:
	echo ./cmd/proxy/main.go | xargs -Ip tinygo build -o p.wasm -target=wasi p
	envoy -c ./cmd/proxy/envoy.yaml --concurrency 2 --log-format '%v'

.PHONY: lint
lint:
	@go run $(golangci_lint) run --build-tags proxytest

.PHONY: format
format:
	@find . -type f -name '*.go' | xargs gofmt -s -w
	@for f in `find . -name '*.go'`; do \
	    awk '/^import \($$/,/^\)$$/{if($$0=="")next}{print}' $$f > /tmp/fmt; \
	    mv /tmp/fmt $$f; \
	done
	@go run $(goimports) -w -local github.com/kzscaler/kzscaler `find . -name '*.go'`

.PHONY: check
check:
	@$(MAKE) format
	@go mod tidy
	@if [ ! -z "`git status -s`" ]; then \
		echo "The following differences will fail CI until committed:"; \
		git diff --exit-code; \
	fi
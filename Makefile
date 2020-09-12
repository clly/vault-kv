.DEFAULT_GOAL := help
.PHONY: help

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	| sed -n 's/^\(.*\): \(.*\)##\(.*\)/\1\3/p' \
	| column -t  -s ' '

clean: ## clean
	rm -f tools/golangci-lint

tools/base: export GOFLAGS = -mod=readonly
tools/base: tools/go.mod tools/go.sum

tools/golangci-lint: tools/base
	cd tools && go build github.com/golangci/golangci-lint/cmd/golangci-lint

tools/go-opine: tools/base
	cd tools && go build oss.indeed.com/go/go-opine

tools: tools/golangci-lint tools/go-opine ## tools

lint: tools/golangci-lint ## lint
	./tools/golangci-lint run 

cover: tools/go-opine ## cover
	./tools/go-opine test -coverprofile=cover.out -min-coverage=0 && go tool cover -html=cover.out && rm -v cover.out

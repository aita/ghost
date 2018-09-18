.DEFAULT_GOAL := build

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build dep dep-update clean

build: ghost ghost-shell  ## Generate an executable file from source files

ghost:
	go build -o ghost cmd/ghost/main.go

ghost-shell:
	go build -o ghost-shell cmd/ghost-shell/main.go

dep:  ## Install dependencies
	dep ensure -v

dep-update:  ## Update dependencies
	dep ensure -v -update

clean:  ## Remove the executable file
	rm -rf ghost ghost-shell

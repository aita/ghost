.DEFAULT_GOAL := help
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build dep dep-update

build:  ## Generate an executable file from source files
	go build -o ghost cmd/ghost/main.go

dep:  ## Install dependencies
	dep ensure -v

dep-update:  ## Update dependencies
	dep ensure -v -update
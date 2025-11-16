.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

.PHONY: check-env
check-env:
ifndef HARKENER_DEPLOY_HOST
	$(error HARKENER_DEPLOY_HOST is undefined)
endif
ifndef HARKENER_DEPLOY_PATH
	$(error HARKENER_DEPLOY_PATH is undefined)
endif

.PHONY: build
build: ## build the binary
	@$(eval os ?= linux)
	@$(eval arch ?= amd64)
	@docker build --build-arg ARCH=$(arch) --build-arg OS=$(os) --output=. . $(args)

.PHONY: deploy
deploy: check-env build ## build the binary and deploy to the server
	@scp harkener-linux-amd64 $(HARKENER_DEPLOY_HOST):$(HARKENER_DEPLOY_PATH)

.PHONY: fmt
fmt: ## reformat
	@go fmt ./...


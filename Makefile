.PHONY: help
help: ## This help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.DEFAULT_GOAL := help

.PHONY: check-env
check-env:
ifndef LISTENER_DEPLOY_HOST
	$(error LISTENER_DEPLOY_HOST is undefined)
endif
ifndef LISTENER_DEPLOY_PATH
	$(error LISTENER_DEPLOY_PATH is undefined)
endif

.PHONY: build
build: ## build the binary
	@docker build --output=. . $(args)

.PHONY: deploy
deploy: check-env build ## build the binary and deploy to the server
	@scp listener $(LISTENER_DEPLOY_HOST):$(LISTENER_DEPLOY_PATH)

.PHONY: fmt
fmt: ## reformat
	@go fmt ./...


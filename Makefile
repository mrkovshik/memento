# Define variables for tools
GOLANGCI_LINT := $(shell which golangci-lint)
PROTOC := $(shell which protoc)

.PHONY: lint
lint: 
	$(GOLANGCI_LINT) run -v --out-format code-climate:gl-code-quality-report.json,line-number --fix

.PHONY: proto
proto:
	$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/memento.proto

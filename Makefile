.PHONY: lint
lint: golangci-lint
	$(GOLANGCI_LINT) run -v --out-format code-climate:gl-code-quality-report.json,line-number --fix


.PHONY: proto
proto: protoc
	--go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/memento.proto
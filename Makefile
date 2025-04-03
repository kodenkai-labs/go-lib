.PHONY: test
test:
	@echo "  >  Running unit tests"
	go clean -testcache
	go test -coverprofile=coverage -cover -v ./internal... ./pkg...

.PHONY: tests-coverage
test-coverage: test
	go tool cover -html=coverage -o cover_out.html
	open cover_out.html

.PHONY: lint
lint: go-lint-install go-lint

.PHONY: go-lint-install
go-lint-install:
ifeq (,$(shell which golangci-lint))
	@echo "  >  Installing golint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- v1.64.5
endif

.PHONY: go-lint
go-lint:
	@echo "  >  Running golint"
	golangci-lint run ./...

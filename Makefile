.PHONY: fmt lint test tidy precommit

fmt:
	go run mvdan.cc/gofumpt@v0.10.0 -w .

tidy:
	go mod tidy

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run

test:
	go test ./...

precommit:
	@if command -v pre-commit >/dev/null 2>&1; then \
		pre-commit run --all-files; \
	else \
		.githooks/pre-commit; \
	fi

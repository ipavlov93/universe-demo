# Development
# Prerequisites: Go

# run tests
run_tests:
	go test ./... -count 1

# formatting
gofmt:
	go fmt ./...

# goimports groups and sorts import sections
goimports:
	goimports --local github.com/ipavlov93/universe-demo/product-sv/ -l -w .

# Generate mocks using mockery
# Prerequisites: mockery.
# Check the latest tool's version on release page.
mockery_install:
	go install github.com/vektra/mockery/v3@v3.3.2

mockery_generate_mocks:
	mockery --config ./internal/mocks/.mockery.yaml

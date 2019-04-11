.PHONY: test
test:
	go test ./... -race -v
lint-prepare:
	@echo "Installing golangci-lint"
	# @go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run \
		--exclude-use-default=false \
		--enable=golint \
		--enable=gocyclo \
		--enable=goconst \
		--enable=unconvert \
		./...
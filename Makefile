BINARY_NAME=org-agenda-cli

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(BINARY_NAME) main.go

.PHONY: test
test:
	go test ./...

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: lint
lint:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: check-fmt
check-fmt:
	test -z "$$(gofmt -l .)"

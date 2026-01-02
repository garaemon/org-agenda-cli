BINARY_NAME=org-agenda-cli

.PHONY: all
all: build

.PHONY: build
build:
	go build -o $(BINARY_NAME) main.go

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: lint
lint:
	go vet ./...

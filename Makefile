GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
BUILD_DIR=build
EXECUTABLE=arbitragebot

help:
	@echo "This is a helper makefile for arbitragebot"
	@echo "Targets:"
	@echo "    lint:	run lint"
	@echo "    check:	run tests"
	@echo "    tidy:	tidy go mod"
	@echo "    build:	builds executable"

$(GOBIN)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.60.1

.PHONY: tools
tools: $(GOBIN)/golangci-lint

lint: tools
	$(GOBIN)/golangci-lint run ./...

check:
	go test ./...

tidy:
	@echo "tidy..."
	go mod tidy

build: clean
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(EXECUTABLE) ./cmd/arbitragebot

clean:
	rm -rf $(BUILD_DIR)
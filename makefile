APP_NAME   := grab
SOURCE     := grab.go
BIN_DIR    := bin
BIN        := $(BIN_DIR)/$(APP_NAME)
#SOURCES    := $(wildcard *.go)

GOFLAGS    := 

GREEN      := \033[1;32m
RED        := \033[1;31m
RESET      := \033[0m

.PHONY: build run test vet fmt tidy clean help
.DEFAULT_GOAL := help

$(BIN): $(SOURCE)
	@mkdir -p $(BIN_DIR)
	@START=$$(date +%s.%3N);                                                               				  \
	if go build $(GOFLAGS) -o $@ $<; then			                                       			  	  \
		END=$$(date +%s.%3N);                                                              				  \
		printf "Compilation $(GREEN)finished$(RESET) in %.2f seconds.\n" $$(echo "$$END - $$START" | bc); \
	else 																				  				  \
		END=$$(date +%s.%3N);                                                             				  \
		printf "Compilation $(RED)exited$(RESET) in %.2f seconds.\n" $$(echo "$$END - $$START" | bc);	  \
		exit 1;																	  				  		  \
	fi

build: $(BIN) ## build the binary

run: build ## compile and run the binary
	@echo "Running $(APP_NAME)..."
	@$(BIN)

test: ## run tests
	@echo "Running tests..."
	@go test ./...

vet: ## run go vet
	@echo "Running go vet..."
	@go vet ./...

fmt: ## run go fmt
	@echo "Running go fmt..."
	@go fmt ./...

tidy: ## run go mod tidy
	@echo "Running go mod tidy..."
	@go mod tidy

clean: ## remove the binary files
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)

help: ## show this help message
	@echo "Makefile for $(APP_NAME)"
	@echo
	@echo "Usage:"
	@echo "\tmake [target]"
	@echo
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "\t%-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo


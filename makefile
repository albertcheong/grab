APP_NAME   := grab
SOURCE     := grab.go
BIN_DIR    := bin
BIN        := $(BIN_DIR)/$(APP_NAME)
#SOURCES    := $(wildcard *.go)

GOFLAGS    := 

VERSION := $(shell git describe --tags --always)
COMMIT  := $(shell git rev-parse --short HEAD)
DATE    := $(shell date +%F)

LDFLAGS    := -X main.Version=$(VERSION) \
              -X main.Commit=$(COMMIT)   \
			  -X main.Date=$(DATE)       \

GREEN      := \033[1;32m
RED        := \033[1;31m
CYAN       := \033[1;36m
PURPLE     := \033[1;35m
RESET      := \033[0m

.PHONY: help
help: ## show this help message
	@echo "$(PURPLE)Usage:$(RESET)"
	@echo "\tmake $(CYAN)[TARGETS]$(RESET)"
	@echo
	@echo "$(PURPLE)Targets:$(RESET)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "\t$(CYAN)%-14s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo

$(BIN): $(SOURCE)
	@mkdir -p $(BIN_DIR)
	@START=$$(date +%s.%3N);                                                               				  \
	if go build -ldflags "$(LDFLAGS)" $(GOFLAGS) -o $@ $<; then			                                  \
		END=$$(date +%s.%3N);                                                              				  \
		printf "Compilation $(GREEN)finished$(RESET) in %.2f seconds.\n" $$(echo "$$END - $$START" | bc); \
	else 																				  				  \
		END=$$(date +%s.%3N);                                                             				  \
		printf "Compilation $(RED)exited$(RESET) in %.2f seconds.\n" $$(echo "$$END - $$START" | bc);	  \
		exit 1;																	  				  		  \
	fi

.PHONY: build
build: $(BIN) ## build the binary

.PHONY: run
run: build ## compile and run the binary
	@$(BIN) $(ARGS)

.PHONY: test
test: ## run tests
	@go test ./...

.PHONY: clean
clean: ## remove the binary files
	@rm -rf $(BIN_DIR)
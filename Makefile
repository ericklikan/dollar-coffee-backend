# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

BUILD_DIR=build
BINARY_NAME=$(BUILD_DIR)/run
    
all: test build

.PHONY: build
build: clean
	$(GOBUILD) -o $(BINARY_NAME) -v

.PHONY: test
test: 
	$(GOTEST) -v ./...

.PHONY: clean
clean: 
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

.PHONY: run
run:
	$(GORUN) main.go

.PHONY: deps
deps:
	$(GOGET) ./...
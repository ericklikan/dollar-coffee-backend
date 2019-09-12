# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run

BUILD_DIR=build_dir
BINARY_NAME=$(BUILD_DIR)/run
    
all: test build
build: clean
	$(GOBUILD) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
run:
	$(GORUN) main.go
deps:
	$(GOGET) ./...
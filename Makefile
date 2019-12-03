# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
MAIN = resolver/cmd/resolver
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=build/resolver

all: test build
build: 
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN) 
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN) ./...
	./$(BINARY_NAME)

.PHONY: all test build run clean

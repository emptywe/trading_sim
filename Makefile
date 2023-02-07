REPONAME := trading_sim

GOCMD=go
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

run:
	 $(GORUN) cmd/simulator/*.go
run-parser:
	 $(GORUN) cmd/parser/*.go

build:
	@cd ./cmd/parser && $(GOBUILD)
	@cd ./cmd/simulator && $(GOBUILD)
.PHONY:
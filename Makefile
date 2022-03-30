GO = GODEBUG=sbrk=1 GO15VENDOREXPERIMENT=1 go
GOFLAGS = -tags netgo -ldflags "-X main.version=$(shell git describe --all --always --long | sed -E 's/^(heads|tags)\///') -X main.buildTime=$(shell date -u "+%s")"
BINARIES = ./bin
CWD = "$(shell pwd)"

build:
	@mkdir -p $(BINARIES)
	@echo "Building Current OS Version"
	$(GO) build $(GOFLAGS) -o $(BINARIES) ./cmd/*

clean: cleanbinary

cleanbinary:
	@echo "Clean binaries folder"
	rm -rf $(BINARIES)/*

.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

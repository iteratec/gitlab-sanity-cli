GOCMD=go
GOBUILD=$(GOCMD) build -ldflags "-extldflags '-static'"
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
MAINGO=cmd/main.go
ARTEFACTSDIR=./

all: setmode clean test build

.PHONY: setmode
setmode:
#	$(GOCMD) env -w GO111MODULE=off

build: freebsd darwin linux windows

freebsd: test
	@echo "Build FreeBSD AMD64"
	env GOOS=freebsd GOARCH=amd64 $(GOBUILD) -o $(ARTEFACTSDIR)/gitlab-sanity-cli.freebsd $(MAINGO)

linux: test
	@echo "Build Linux AMD64"
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(ARTEFACTSDIR)/gitlab-sanity-cli.linux $(MAINGO)

darwin: test
	@echo "Build MacOS AMD64"
	env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(ARTEFACTSDIR)/gitlab-sanity-cli.darwin $(MAINGO)

windows: test
	@echo "Build Windows AMD64"
	env GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(ARTEFACTSDIR)/gitlab-sanity-cli.windows $(MAINGO)

# dbg-build:
# 	$(GOBUILD) -v -gcflags=all="-N -l" -tags debug $(MAINFOLDER)

test:
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)

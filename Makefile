ARTEFACTSDIR=.
GOCMD=go
GOBUILD=$(GOCMD) build -ldflags "-extldflags '-static'"
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
MAINGO=cmd/main.go
BINARYNAME=gitlab-sanity-cli
TAR=tar czvf
ZIP=zip --junk-paths
SHASUM=shasum

all: setmode clean test build

.PHONY: setmode
setmode:
#	$(GOCMD) env -w GO111MODULE=off

build: freebsd darwin linux windows

freebsd: test
	@echo "Build FreeBSD AMD64"
	env GOOS=freebsd GOARCH=amd64 $(GOBUILD) -o $(ARTEFACTSDIR)/$(BINARYNAME).freebsd $(MAINGO)

linux: test
	@echo "Build Linux AMD64"
	env GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(ARTEFACTSDIR)/$(BINARYNAME).linux $(MAINGO)

darwin: test
	@echo "Build MacOS AMD64"
	env GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(ARTEFACTSDIR)/$(BINARYNAME).darwin $(MAINGO)

windows: test
	@echo "Build Windows AMD64"
	env GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(ARTEFACTSDIR)/$(BINARYNAME).windows $(MAINGO)

archive: 
	$(TAR) $(ARTEFACTSDIR)/$(BINARYNAME).freebsd.amd64.tar.gz $(ARTEFACTSDIR)/$(BINARYNAME).freebsd
	$(TAR) $(ARTEFACTSDIR)/$(BINARYNAME).linux.amd64.tar.gz $(ARTEFACTSDIR)/$(BINARYNAME).linux
	$(TAR) $(ARTEFACTSDIR)/$(BINARYNAME).darwin.amd64.tar.gz $(ARTEFACTSDIR)/$(BINARYNAME).darwin
	$(ZIP) $(ARTEFACTSDIR)/$(BINARYNAME).windows.zip $(ARTEFACTSDIR)/$(BINARYNAME).windows

hash:
	$(SHASUM) $(ARTEFACTSDIR)/$(BINARYNAME).freebsd.amd64.tar.gz > $(ARTEFACTSDIR)/$(BINARYNAME).freebsd.amd64.tar.gz.sha256
	$(SHASUM) $(ARTEFACTSDIR)/$(BINARYNAME).linux.amd64.tar.gz > $(ARTEFACTSDIR)/$(BINARYNAME).linux.amd64.tar.gz.sha256
	$(SHASUM) $(ARTEFACTSDIR)/$(BINARYNAME).darwin.amd64.tar.gz > $(ARTEFACTSDIR)/$(BINARYNAME).darwin.amd64.tar.gz.sha256
	$(SHASUM) $(ARTEFACTSDIR)/$(BINARYNAME).windows.zip > $(ARTEFACTSDIR)/$(BINARYNAME).windows.zip.sha256

release: archive hash

# dbg-build:
# 	$(GOBUILD) -v -gcflags=all="-N -l" -tags debug $(MAINFOLDER)

test:
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)
	@find . -name "gitlab-sanity-cli.*" -type f -delete
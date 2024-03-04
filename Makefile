NAME = docker2exe
OUTPUT = dist
VERSION = v0.3.0
SOURCES = $(wildcard *.go)
TARGETS = darwin/amd64 darwin/arm64 linux/amd64 windows/amd64

OS = $(shell go env GOOS)
ARCH = $(shell go env GOARCH)

os = $(word 1, $(subst /, ,$@))
arch = $(word 2, $(subst /, ,$@))

.PHONY: all
all: $(TARGETS)

.PHONY: clean
clean:
	$(RM) -rf $(OUTPUT)

.PHONY: test
test: all
	dist/docker2exe-$(OS)-$(ARCH) --name test --image alpine
	dist/test-$(OS)-$(ARCH) echo OK
	dist/docker2exe-$(OS)-$(ARCH) --name test-embed --image alpine --embed
	dist/test-embed-$(OS)-$(ARCH) echo OK

$(OUTPUT):
	mkdir $(OUTPUT)

$(TARGETS): $(SOURCES) $(OUTPUT)
	GOOS=$(os) GOARCH=$(arch) go build -o "$(OUTPUT)/$(NAME)-$(os)-$(arch)"

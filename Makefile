REPO_URI ?= github.com/stackanetes
REPO_PATH ?= $(REPO_URI)/kubernetes-entrypoint

prepare:
	@echo "Preapre GOPATH"
	test -h gopath/src/$(REPO_PATH) || \
		( mkdir -p gopath/src/$(REPO_URI); \
		ln -s ../../../.. gopath/src/$(REPO_PATH) )

build: prepare
	@echo "Building kubernetes-entrypoint for $(GOOS)/$(GOARCH) $(GOPATH)"
	mkdir -p bin/$(GOARCH)
	go build -o bin/$(GOARCH)/kubernetes_entrypoint

linux-arm64:
	export GOOS="linux"; \
	export GOARCH="arm64"; \
	export GOPATH="$(PWD)/gopath"; \
	$(MAKE) build

linux-amd64:
	export GOOS="linux"; \
	export GOARCH="amd64"; \
	export GOPATH="$(PWD)/gopath"; \
	$(MAKE) build

clean:
	rm -rf gopath
	rm -rf bin

test: prepare
	export GOPATH="$(PWD)/gopath"; \
	go test

all: linux-amd64 linux-arm64

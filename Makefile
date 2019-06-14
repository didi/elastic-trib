default: help

HOST_GOLANG_VERSION := $(shell go version | cut -d ' ' -f3 | cut -c 3-)
MODULE := elastic-trib
ifneq (,$(wildcard .git/.*))
    COMMIT = $(shell git rev-parse HEAD 2> /dev/null || true)
    VERSION	= $(shell git describe --tags --abbrev=0 2> /dev/null)
else
    COMMIT = "unknown"
    VERSION = "unknown"
endif

export GOPATH := $(shell cd ./ && pwd)/vendor
export CURDIR := $(shell cd ./ && pwd)

## Make bin for $MODULE.
bin: 
	@echo "GOVERSION: ${HOST_GOLANG_VERSION}"
	@echo "GOPATH:" $$GOPATH
	@ln  -s "${CURDIR}/vendor" "${CURDIR}/vendor/src"
	go build -i -ldflags "-X main.gitCommit=${COMMIT} -X main.version=${VERSION}" -o ${MODULE} .; \
	    rm -rf ${CURDIR}/vendor/src;  rm -rf ${CURDIR}/vendor/pkg

## Build debug trace for $MODULE.
debug:
	@echo "GOVERSION: ${HOST_GOLANG_VERSION}"
	@echo "GOPATH:" $$GOPATH
	@ln  -s "${CURDIR}/vendor" "${CURDIR}/vendor/src"
	go build -n -v -i -ldflags "-X main.gitCommit=${COMMIT} -X main.version=${VERSION}" -o ${MODULE} .; \
	    rm -rf ${CURDIR}/vendor/src

## Make static link bin for $MODULE.
static-bin:
	@echo "GOVERSION:" ${HOST_GOLANG_VERSION}
	@echo "GOPATH:" $$GOPATH
	@ln  -s "${CURDIR}/vendor" "${CURDIR}/vendor/src"
	go build -i -ldflags "-w -extldflags -static -X main.gitCommit=${COMMIT} -X main.version=${VERSION}" -o ${MODULE} .; \
	    rm -rf ${CURDIR}/vendor/src

## Get dep tool for managing dependencies for Go projects.
dep:
	go get -u github.com/golang/dep/cmd/dep

## Get dep tool and init project. 
depinit: dep
	dep init
	dep ensure

## Get vet go tools.
vet:
	go get golang.org/x/tools/cmd/vet

## Validate this go project.
validate:
	script/validate-gofmt
	#go vet ./...

## Run test case for this go project.
test:
	go test -v ./...

## Clean everything (including stray volumes).
clean:
#	find . -name '*.created' -exec rm -f {} +
	-rm -rf ${CURDIR}/vendor/src
	-rm -rf ${CURDIR}/vendor/pkg
	-rm -rf ${CURDIR}/vendor/vendor
	-rm -rf var
	-rm -f ${MODULE}

help: # Some kind of magic from https://gist.github.com/rcmachado/af3db315e31383502660
	$(info Available targets)
	@awk '/^[a-zA-Z\-\_0-9]+:/ {                                   \
		nb = sub( /^## /, "", helpMsg );                             \
		if(nb == 0) {                                                \
			helpMsg = $$0;                                             \
			nb = sub( /^[^:]*:.* ## /, "", helpMsg );                  \
		}                                                            \
		if (nb) {                                                     \
			h = sub( /[^ ]*MODULE/, "'${MODULE}'", helpMsg );        \
			printf "   \033[1;31m%-" width "s\033[0m %s\n", $$1, helpMsg; \
		}															\
	}                                                              \
	{ helpMsg = $$0 }'                                             \
	width=$$(grep -o '^[a-zA-Z_0-9]\+:' $(MAKEFILE_LIST) | wc -L 2> /dev/null)  \
	$(MAKEFILE_LIST)


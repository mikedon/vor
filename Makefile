# Check for required command tools to build or stop immediately
EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
	$(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

ROOT_DIR     := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BINARY        = vor
VERSION       = 0.1.0
BUILD         = `git rev-parse HEAD`
PLATFORMS     = darwin linux windows
ARCHITECTURES = 386 amd64

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

default: build

all: clean build_all zip #install

build:
	go build ${LDFLAGS} -o ${BINARY}

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o $(BINARY)-$(GOOS)-$(GOARCH))))

zip:
	$(eval LOCAL_BINARIES := $(shell find * -type f -maxdepth 1 -name "vor-*-*"))
	tar -zcf ${BINARY}-${VERSION}.tar.gz $(LOCAL_BINARIES) README.md LICENSE && \
	zip -r ${BINARY}-${VERSION}.zip $(LOCAL_BINARIES) README.md LICENSE && \
	echo `md5 ${BINARY}-${VERSION}.zip ${BINARY}-${VERSION}.tar.gz`

install:
	go install ${LDFLAGS}

# Remove only what we've created
clean:
	find * -name '${BINARY}-[a-zA-Z0-9]*-[a-zA-Z0-9]*' -delete && \
	find * -name '${BINARY}-[0-9]*.zip' -delete && \
	find * -name '${BINARY}-[0-9]*.tar.gz' -delete

.PHONY: check clean install build_all all

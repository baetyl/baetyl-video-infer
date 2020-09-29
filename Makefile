MODULE   := video-infer
BIN      := baetyl-$(MODULE)
OUTPUT   := output
SRC_FILES:= $(shell find . -type f -name '*.go')

GIT_TAG := $(shell git tag --contains HEAD)
GIT_REV := git-$(shell git rev-parse --short HEAD)
VERSION := $(if $(GIT_TAG),$(GIT_TAG),$(GIT_REV))

GO       := go
GO_MOD   := $(GO) mod
GO_ENV   := env GO111MODULE=on GOPROXY=https://goproxy.cn CGO_ENABLED=1
GO_TAGS  :=
GO_FLAGS := -ldflags '-X "github.com/baetyl/baetyl-go/v2/utils.REVISION=$(GIT_REV)" -X "github.com/baetyl/baetyl-go/v2/utils.VERSION=$(VERSION)"'
GO_BUILD := $(GO_ENV) $(GO) build $(GO_TAGS) $(GO_FLAGS)

XFLAGS     := --push
XPLATFORMS := linux/amd64,linux/arm64,linux/arm/v7
REGISTRY   := baetyl

.PHONY: all
all: build test

.PHONY: build
build: $(SRC_FILES)
	$(GO_BUILD) -o $(OUTPUT)/$(BIN) .

.PHONY: image
image: clean
	docker build -t $(REGISTRY)/$(MODULE):$(VERSION) -f Dockerfile .

.PHONY: image-openvino
image-openvino: clean
	docker build -t $(REGISTRY)/$(MODULE)-openvino:$(VERSION) -f openvino/Dockerfile .

.PHONY: image-all
image-all:
	@-docker buildx create --name baetyl
	@docker buildx use baetyl
	@docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	docker buildx build $(XFLAGS) --platform $(XPLATFORMS) -t $(REGISTRY)/$(MODULE):$(VERSION) -f buildx/Dockerfile .

.PHONY: fmt
fmt:
	$(GO_MOD) tidy
	@go fmt ./...

.PHONY: clean
clean:
	@rm -rf $(OUTPUT) $(BIN)


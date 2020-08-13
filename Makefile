PREFIX?=/usr/local
VERSION?=latest
SRC=$(wildcard *.go)

all: build

build: $(SRC)
	@echo "BUILD $@"
	@go build ${GO_FLAGS} .

image: build
	@echo "BUILDX $<"
	@docker build -t $(DOCKER_REPO)$<$(IMAGE_SUFFIX):$(VERSION) .

.PHONY: clean
clean:
	rm -f baetyl-video-infer package.zip

.PHONY: rebuild
rebuild: clean all

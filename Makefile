MODULE:=baetyl-video-infer
BUILD_TAGS:=

all: build

build:
	env GO111MODULE=on GOPROXY=https://goproxy.cn go build -tags $(BUILD_TAGS) -o $(MODULE) .

image: build
	docker build -t $(MODULE) .

.PHONY: clean
clean:
	rm -f $(MODULE)

.PHONY: rebuild
rebuild: clean build

.PHONY: fmt
fmt:
	go fmt ./...
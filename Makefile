MODULE:=baetyl-video-infer

all: build

build:
	@env GO111MODULE=on GOPROXY=https://goproxy.cn go build -o $(MODULE) .

image: build
	@docker build -t $(MODULE) .

.PHONY: clean
clean:
	rm -f $(MODULE)

.PHONY: rebuild
rebuild: clean build

.PHONY: fmt
fmt:
	go fmt ./...
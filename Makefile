VERSION?=0.4.7

all: build
build:
	CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -v .

docker:
	docker build -t pool:0.9 .


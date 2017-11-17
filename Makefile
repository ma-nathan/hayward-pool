VERSION?=0.4.4

all: build
build:
	go build -ldflags "-X main.Version=$(VERSION)" -v .


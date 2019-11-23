VERSION?=0.4.6

all: build
build:
	go build -ldflags "-X main.Version=$(VERSION)" -v .


VERSION?=0.4.5

all: build
build:
	go build -ldflags "-X main.Version=$(VERSION)" -v .


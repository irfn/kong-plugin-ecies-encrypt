.PHONY: all

all: build build-plugin

build:
	go build -o ecies-encrypt ecies-encrypt.go

build-plugin:
	go build -buildmode=plugin ecies-encrypt.go -o ecies-encrypt.so
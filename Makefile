BINARY=bin/cep

.PHONY: all
all: build

.PHONY: build
build: dependencies build-bin

build-bin:
	go build -o ./$(BINARY)

clean:
	rm -rf bin/*

dependencies:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: run
run: build-bin
	./$(BINARY)

.PHONY: test
test:
	go test ./... -race -cover


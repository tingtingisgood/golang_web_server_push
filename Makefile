GO=/usr/local/go/bin/go

all: golangWebServerPush

.PHONY: clean lint

SRC=$(shell find . -maxdepth 2 -name "*.go" -type f)


lint:
	golint
	golint monitor

golangWebServerPush: $(SRC)
	$(GO) build ./monitor
	$(GO) build -o $@

clean:
	rm golangWebServerPush


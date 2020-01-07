all: run
.PHONY: all

mod:
	go mod vendor
	go mod tidy
.PHONY: mod

build: mod
	go build
.PHONY: build

test:
	go test
.PHONY: test

run: build
	./kobj
.PHONY: run

clean:
	rm -f ./kobj
	rm -rf ./vendor/
.PHONY: clean

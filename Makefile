# GO111MODULE ?= on
# GOFLAGS 	?= -mod=vendor

all: run
	@echo "done."
.PHONY: all

run: build
	./kobj server
.PHONY: run

build: mod
	go build
.PHONY: build

mod:
	go mod tidy
	go mod vendor
.PHONY: mod

gen: mod
	bash ./vendor/k8s.io/code-generator/generate-internal-groups.sh \
		deepcopy,openapi,client,protobuf \
		github.com/kobj-io/kobj/pkg/apis \
		github.com/kobj-io/kobj/pkg/apis \
		github.com/kobj-io/kobj/pkg/apis \
		"kobj:v1alpha1" \
		--go-header-file ./boilerplate.go.txt
.PHONY: gen

test:
	go test ./...
.PHONY: test

clean:
	rm -f ./kobj
	rm -rf ./vendor/
.PHONY: clean

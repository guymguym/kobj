# GO111MODULE ?= on
# GOFLAGS 	?= -mod=vendor

all: run
.PHONY: all

mod:
	go mod vendor
	go mod tidy
.PHONY: mod

gen: mod
	bash ./vendor/k8s.io/code-generator/generate-internal-groups.sh \
		deepcopy,openapi \
		github.com/kobj-io/kobj/pkg/apis \
		github.com/kobj-io/kobj/pkg/apis \
		github.com/kobj-io/kobj/pkg/apis \
		"kobj:v1alpha1" \
		--go-header-file ./boilerplate.go.txt
.PHONY: gen

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

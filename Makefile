# export GO111MODULE = on
# export GOFLAGS = -mod=vendor

all: build
	@echo "✅ all: done."
.PHONY: all

build: mod
	go build
	@echo "✅ build: done."
.PHONY: build

mod:
	go mod tidy
	go mod vendor
.PHONY: mod

test:
	go test ./...
	@echo "✅ test: done."
.PHONY: test

clean:
	rm -f ./kobj
	rm -rf ./vendor/
	@echo "✅ clean: done."
.PHONY: clean

#-------------------#
#- Code Generation -#
#-------------------#

gen: gen-groups gen-internal-groups # gen-proto
	@echo "✅ gen: done."
.PHONY: gen

gen-groups:
	bash ./vendor/k8s.io/code-generator/generate-groups.sh \
		all \
		github.com/kobj-io/kobj/pkg/generated \
		github.com/kobj-io/kobj/pkg/apis \
		"kobj:v1alpha1" \
		--go-header-file ./boilerplate.go.txt
.PHONY: gen-groups

gen-internal-groups:
	bash ./vendor/k8s.io/code-generator/generate-internal-groups.sh \
		deepcopy,openapi \
		github.com/kobj-io/kobj/pkg/generated \
		github.com/kobj-io/kobj/pkg/apis \
		github.com/kobj-io/kobj/pkg/apis \
		"kobj:v1alpha1" \
		--go-header-file ./boilerplate.go.txt
.PHONY: gen-internal-groups

gen-proto:
	go run ./vendor/k8s.io/code-generator/cmd/go-to-protobuf \
		-p github.com/kobj-io/kobj/pkg/apis \
		-o ./pkg/generated/ \
		--proto-import $(PWD)/vendor/k8s.io/apimachinery/pkg/util/intstr/,$(PWD)/vendor/github.com/gogo/protobuf/gogoproto/
.PHONY: gen-proto

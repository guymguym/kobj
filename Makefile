# export GO111MODULE = on
# export GOFLAGS = -mod=vendor

all: image
	@echo "✅ all: done."
.PHONY: all

image: mod
	GOOS=linux GOARCH=amd64 go build -o kobj-linux
	docker build -t kobj/kobj .
.PHONY: image

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

help:
	@echo '#--------#'
	@echo '#- Help -#'
	@echo '#--------#'
	@echo
	@echo '# Run Server:'
	@echo
	@echo '    ./kobj server'
	@echo
	@echo '# Run kobj client:'
	@echo
	@echo '    ./kobj put -n default aaa < main.go'
	@echo '    ./kobj get -n default aaa'
	@echo
	@echo '# Run kubectl:'
	@echo
	@echo '    export KUBECONFIG=./kubeconfig-local'
	@echo '    kubectl get kobjs.kobj.io aaa'
	@echo '    kubectl get kobjs.kobj.io aaa -o yaml'
	@echo
	@echo '# Run curl:'
	@echo
	@echo '    curl -k https://localhost:8443'
	@echo '    curl --cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt' \
		'-H "Authorization: Bearer $$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)"' \
		'https://kubernetes.default.svc/'
	@echo
.PHONY: help

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

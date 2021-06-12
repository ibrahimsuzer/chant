PROTO_ROOT=proto
PROTO_PATHS=

asdf:
	asdf plugin add golang; if [ $$? -eq 2 ] ; then true ; fi
	asdf plugin add protoc; if [ $$? -eq 2 ] ; then true ; fi
	asdf plugin add goreleaser; if [ $$? -eq 2 ] ; then true ; fi
	asdf plugin add golangci-lint; if [ $$? -eq 2 ] ; then true ; fi
	asdf install

install:
	go install github.com/bufbuild/buf/cmd/buf@v0.43.2
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking@v0.43.2
	go install github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@v0.43.2
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.4.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.4.0
	go install github.com/jfeliu007/goplantuml/cmd/goplantuml@v1.5.2
	go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v1.4.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26.0
	go install github.com/grpc-ecosystem/grpc-health-probe@v0.4.2

lint:
	buf lint

check:
	buf breaking --against .git#branch=main

build:
	buf build

generate:
ifdef PROTO_PATHS
	buf generate --path ${PROTO_PATHS}
endif
	goplantuml -recursive \
		-show-aggregations \
		-aggregate-private-members \
		-show-connection-labels \
		-show-aliases \
		-show-compositions \
		-show-implementations \
		-show-options-as-note \
		-ignore $(shell find * -type f \( -path "*.pb.go" -o -path "*_mock.go" \) -printf '%h\n' | sort | uniq | paste -sd, -) -output diagram.puml .

#!/bin/bash

# set -eux
set -e
set -o pipefail

cd "$(git rev-parse --show-toplevel)"

# generate *.pb.go or wire_gen.go
go generate ./...

third_party="./third_party"
grpc_gateway="./third_party/grpc-gateway"
core_outdir="./internal/corepb"
core_indir="./protobuf/core"
swagger_outdir="./cmd/flyio/controller/doc/static/swagger"
# remove *.go
rm -fv $core_outdir/*.go
rm -fv $swagger_outdir/*.json

for item in "$core_indir"/* ; do
    echo "$item"
    if [[ -f "$item" ]]; then
        protoc \
        --go_out="$core_outdir" \
        --go-grpc_out="$core_outdir" \
        --go_opt=paths=source_relative \
        --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out="$core_outdir" \
        --grpc-gateway_opt logtostderr=true \
        --grpc-gateway_opt paths=source_relative \
        --grpc-gateway_opt generate_unbound_methods=true \
        --grpc-gateway_opt register_func_suffix=GW \
        --grpc-gateway_opt allow_delete_body=true \
        --openapiv2_out "$swagger_outdir" \
        --openapiv2_opt logtostderr=true \
        --proto_path="$third_party" \
        --proto_path="$grpc_gateway" \
        --proto_path="$core_indir" \
        "$item"
    fi
done

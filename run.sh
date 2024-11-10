#!/bin/bash
DIR="$( cd "$( dirname "$0" )" && pwd )"

build_lib() {
    if [ ! -n "$BITCOIN_VAULT_DIR" ]; then
        echo "Please set the BITCOIN_VAULT_DIR environment variable"
        exit 1
    fi
    if [ ! -d "$BITCOIN_VAULT_DIR" ]; then
        echo "bitcoin-vault directory not found: $BITCOIN_VAULT_DIR"
        exit 1
    fi
    cd $BITCOIN_VAULT_DIR
    cargo build --release
    mkdir -p $DIR/lib
    cp target/release/libbitcoin_vault_ffi.* $DIR/lib
    export CGO_LDFLAGS="-L./lib -lbitcoin_vault_ffi" 
    export CGO_CFLAGS="-I./lib"
}
btcsigner() {
    docker compose -f docker/docker-compose.yaml up -d bitcoin-testnet4
}
 # example ./run.sh test internals/signer/handlers
test() {
    export CGO_LDFLAGS="-L./lib -lbitcoin_vault_ffi" 
    export CGO_CFLAGS="-I./lib"
    if [ -n "$1" ]; then
        go test -timeout 10m github.com/scalarorg/protocol-signer/${1} -v -count=1
    else
        go test -timeout 10m -v -count=1
    fi
    
}
start() {
    go run ./ start --config example/config.toml
}

$@
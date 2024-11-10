# protocol-signer

Singning service for protocol. This service is built on top of bitcoin-vault.

## Build lib

Set the `BITCOIN_VAULT_DIR` environment variable to the path of the bitcoin-vault directory.
for example:

```
export BITCOIN_VAULT_DIR=$HOME/workspace/codelight/scalar.org/bitcoin-vault
```

Then run the following command to build the lib:

```
./run.sh build_lib
```

## Run bitcoin signer node

```
./run.sh btcsigner
```

## Run tests

```
./run.sh test {test_path}
```

For example:

```
./run.sh test internals/signer/handlers
```

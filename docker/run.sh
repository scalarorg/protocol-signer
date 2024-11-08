#!/bin/sh

createwallet() {
    echo "Creating wallet"
    bitcoin-cli createwallet "protocol" false false protocol false false true
}

importwallet() {
    echo "Importing wallet"
    bitcoin-cli walletpassphrase protocol 60
    bitcoin-cli importprivkey "cVpL6mBRYV3Dmkx87wfbtZ4R3FTD6g58VkTt1ERkqGTMzTcDVw5M"
}

dumpprivkey() {
    echo "Dumping private key"
    bitcoin-cli dumpprivkey $1 >>/root/privkey.txt
}

wait_for_bitcoin_port() {
    echo "Waiting for Bitcoin RPC port 18332..."
    while ! timeout 1 bash -c "echo > /dev/tcp/localhost/18332" 2>/dev/null; do
        echo "Still waiting for Bitcoin RPC port..."
        sleep 1
    done
    echo "Bitcoin RPC port is available"
}

entrypoint() {
    bitcoind $@
    wait_for_bitcoin_port
    createwallet
    importwallet
    dumpprivkey "tb1q37dgjm7e7h385aykhd6gps7uqx0kv26w2ugu8c"
    sleep infinity
}

$@

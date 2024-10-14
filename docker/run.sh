#!/bin/sh
# wget https://bitcoincore.org/bin/bitcoin-core-26.0/bitcoin-26.0-x86_64-linux-gnu.tar.gz -o /root/bitcoin-26.0-x86_64-linux-gnu.tar.gz
# wget https://bitcoincore.org/bin/bitcoin-core-26.0/bitcoin-26.0-x86_64-linux-gnu.tar.gz
# tar -xvf bitcoin-26.0-x86_64-linux-gnu.tar.gz
# sudo ln -sf bitcoin-26.0/bin/bitcoind /usr/bin/bitcoind
# bitcoind
start_bitcoind() {
    bitcoind -testnet \
      -rpcbind=${RPC_BIND:-127.0.0.1:18332} \
      -rpcuser=${RPC_USER:-user} \
      -rpcpassword=${RPC_PASS:-password} \
      -rpcallowip=${RPC_ALLOWIP:-127.0.0.1/0} \
      -datadir=${DATADIR:-/data/.bitcoin} \
      -server=${SERVER:-1} \
      -txindex=${TXINDEX:-1} \
      -connect=${CONNECT:-0} \
      -daemon=${DAEMON:-1}
}
createwallet() {
   bitcoin-cli -named createwallet \
        wallet_name=${WALLET_NAME:-covenant} \
        passphrase=${WALLET_PASSPHRASE:-covenant} \
        load_on_startup=true \
        descriptors=false # create legacy wallet
}
getnewaddress() {
    BTC_ADDRESS=$(bitcoin-cli getnewaddress)
    echo $BTC_ADDRESS>$DATADIR/address.txt
    bitcoin-cli walletpassphrase ${WALLET_PASSPHRASE:-covenant} 60
    bitcoin-cli getaddressinfo $BTC_ADDRESS>$DATADIR/addressinfo.txt
    bitcoin-cli dumpprivkey $BTC_ADDRESS>$DATADIR/privkey.txt
}
entrypoint() {
    bitcoind
    while ! nc -z 127.0.0.1 18332; do
        sleep 1
    done
    createwallet
    getnewaddress
    sleep infinity
}

$@   
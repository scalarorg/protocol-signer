FROM golang:1.23.2-alpine as builder

# Version to build. Default is the Git HEAD.
ARG VERSION="HEAD"

# Use muslc for static libs
ARG BUILD_TAGS="muslc"

RUN apk add --no-cache --update openssh git make build-base linux-headers libc-dev \
    pkgconfig zeromq-dev musl-dev alpine-sdk libsodium-dev \
    libzmq-static libsodium-static gcc

# Build
WORKDIR /go/src/github.com/scalarorg/protocol-signer
# Cache dependencies
COPY go.mod go.sum /go/src/github.com/scalarorg/protocol-signer/

RUN go mod download

# Copy the rest of the files
COPY ./ /go/src/github.com/scalarorg/protocol-signer/

ARG BITCOIN_VAULT_PATH
RUN BITCOIN_VAULT_PATH=$(go list -f '{{.Dir}}' github.com/scalarorg/bitcoin-vault/ffi/go-psbt) && \
    echo "BITCOIN_VAULT_PATH=$BITCOIN_VAULT_PATH" >>/etc/profile && \
    source /etc/profile && \
    export BITCOIN_VAULT_PATH_ENV=$BITCOIN_VAULT_PATH

RUN echo $BITCOIN_VAULT_PATH

RUN CGO_LDFLAGS="$CGO_LDFLAGS -lstdc++ -lm -lsodium -L${BITCOIN_VAULT_PATH}/lib/linux -lbitcoin_vault_ffi" \
    CGO_CFLAGS="-I${BITCOIN_VAULT_PATH}/lib/linux" \
    CGO_ENABLED=1 \
    BUILD_TAGS=$BUILD_TAGS \
    LINK_STATICALLY=true \
    make build

# FINAL IMAGE
# FROM alpine:3.16 AS run

# RUN addgroup --gid 1138 -S protocol-signer && adduser --uid 1138 -S protocol-signer -G protocol-signer

# RUN apk add bash curl jq

# # Copy the binary

# COPY --from=builder /go/src/github.com/scalarorg/protocol-signer/build/protocol-signer /bin/protocol-signer

# # Copy the bitcoin-vault library files
# COPY --from=builder ${BITCOIN_VAULT_PATH}/lib/linux /usr/local/lib/

# RUN ls -l /usr/local/lib

# WORKDIR /home/protocol-signer
# RUN chown -R protocol-signer /home/protocol-signer
# USER protocol-signer

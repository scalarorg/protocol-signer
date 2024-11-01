FROM rust:1.82-alpine3.20 as libbuilder
RUN apk add --no-cache git libc-dev
# Build bitcoin-vault lib
# Todo: select a specific version
RUN git clone https://github.com/scalarorg/bitcoin-vault.git
WORKDIR /bitcoin-vault
RUN cargo build --release

FROM golang:1.23.2-alpine3.20 as builder

# Version to build. Default is the Git HEAD.
ARG VERSION="HEAD"

# Use muslc for static libs
ARG BUILD_TAGS="muslc"

RUN apk add --no-cache --update openssh git make build-base linux-headers libc-dev \
    pkgconfig zeromq-dev musl-dev alpine-sdk libsodium-dev \
    libzmq-static libsodium-static gcc


COPY --from=libbuilder /bitcoin-vault/target/release/libbitcoin_vault_ffi.* /usr/lib/

# Build
WORKDIR /go/src/github.com/scalarorg/protocol-signer
# Cache dependencies
COPY go.mod go.sum /go/src/github.com/scalarorg/protocol-signer/

RUN go mod download

# Copy the rest of the files
COPY ./ /go/src/github.com/scalarorg/protocol-signer/
# RUN ldconfig /usr/lib
RUN CGO_LDFLAGS="$CGO_LDFLAGS -lstdc++ -lm -lsodium -lbitcoin_vault_ffi" \
    CGO_ENABLED=1 \
    BUILD_TAGS=$BUILD_TAGS \
    LINK_STATICALLY=true \
    make build

# FINAL IMAGE
FROM alpine:3.20 AS run

RUN addgroup --gid 1138 -S protocol-signer && adduser --uid 1138 -S protocol-signer -G protocol-signer

RUN apk add bash curl jq

# Copy the binary

COPY --from=builder /go/src/github.com/scalarorg/protocol-signer/build/protocol-signer /bin/protocol-signer

WORKDIR /home/protocol-signer
RUN chown -R protocol-signer /home/protocol-signer
USER protocol-signer

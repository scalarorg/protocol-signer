#!/bin/bash

run() {
    go run ./ start --config example/config.toml
}

$@
#!/bin/bash

run() {
    go run ./ start --config example/config.toml --params example/global-params.json
}

$@
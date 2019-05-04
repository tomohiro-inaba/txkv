#!/bin/bash

function write() {
    sleep $((RANDOM % 10))
    key=$(pwgen --no-capitalize --no-numerals 1 1)
    value=$(echo $RANDOM)
    TXID=$(curl -s 'http://localhost:8080/begin?writable=true' | sed 's/txid=//')
    curl -s "http://localhost:8080/write/${key}?value=${value}&txid=${TXID}"
    curl -s "http://localhost:8080/commit?txid=${TXID}"
}

for i in $(seq 1 100); do
    (write) &
done
wait

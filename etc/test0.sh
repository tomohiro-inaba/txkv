#!/bin/bash

TXID=$(curl -s 'http://localhost:8080/begin?writable=true' | sed 's/txid=//')
curl -s "http://localhost:8080/write/a?value=100&txid=${TXID}"
curl -s "http://localhost:8080/commit?txid=${TXID}"

TXID=$(curl -s 'http://localhost:8080/begin?writable=false' | sed 's/txid=//')
curl -s "http://localhost:8080/read/a?txid=${TXID}"
curl -s "http://localhost:8080/commit?txid=${TXID}"

TXID=$(curl -s 'http://localhost:8080/begin?writable=true' | sed 's/txid=//')
curl -s "http://localhost:8080/write/a?value=200&txid=${TXID}"
curl -s "http://localhost:8080/commit?txid=${TXID}"

TXID=$(curl -s 'http://localhost:8080/begin?writable=false' | sed 's/txid=//')
curl -s "http://localhost:8080/read/a?txid=${TXID}"
curl -s "http://localhost:8080/commit?txid=${TXID}"

TXID=$(curl -s 'http://localhost:8080/begin?writable=true' | sed 's/txid=//')
curl -s "http://localhost:8080/write/b?value=300&txid=${TXID}"
curl -s "http://localhost:8080/commit?txid=${TXID}"

TXID=$(curl -s 'http://localhost:8080/begin?writable=false' | sed 's/txid=//')
curl -s "http://localhost:8080/read/b?txid=${TXID}"
curl -s "http://localhost:8080/commit?txid=${TXID}"

TXID=$(curl -s 'http://localhost:8080/begin?writable=true' | sed 's/txid=//')
curl -s "http://localhost:8080/write/b?value=400&txid=${TXID}"
curl -s "http://localhost:8080/commit?txid=${TXID}"

TXID=$(curl -s 'http://localhost:8080/begin?writable=false' | sed 's/txid=//')
curl -s "http://localhost:8080/read/b?txid=${TXID}"
curl -s "http://localhost:8080/commit?txid=${TXID}"

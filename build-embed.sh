#!/bin/bash

if [[ $# -ne 1 ]]; then
    echo "usage: $0 \$REMOTE_URL"
    exit 92
fi

export REMOTE_URL=$1
# Genrate key and cert 
mkcert --key-file key.pem -cert-file cert.pem $REMOTE_URL 127.0.0.1 ::1

# Embed them in binary
CERT=$(cat cert.pem)
KEY=$(cat key.pem)
CGO_ENABLED=0 go build -ldflags "-X 'github.com/ariary/console.sh/pkg/console.EmbedCert=${CERT}' -X 'github.com/ariary/console.sh/pkg/console.EmbedKey=${KEY}'" cmd/console.sh/console.sh.go

rm cert.pem key.pem
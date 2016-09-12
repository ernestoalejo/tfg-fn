#!/bin/bash

set -eu

if [ ! -f certs/ca.pem ]; then
  echo "--- Generate CA cert"
  cfssl gencert -initca certs/ca-csr.json | cfssljson -bare certs/ca
fi

if [ ! -f certs/server.pem ]; then
  echo "--- Generate server cert"
  cfssl gencert -ca=certs/ca.pem -ca-key=certs/ca-key.pem -config=certs/ca-config.json -profile=server certs/server-csr.json | cfssljson -bare certs/server
fi

if [ ! -f certs/clients/fnctl.pem ]; then
  echo "--- Generate fnctl cert"
  cfssl gencert -ca=certs/ca.pem -ca-key=certs/ca-key.pem -config=certs/ca-config.json -profile=client certs/clients/fnctl-csr.json | cfssljson -bare certs/clients/fnctl
fi

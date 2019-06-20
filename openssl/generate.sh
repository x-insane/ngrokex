#!/bin/bash

# goto openssl direcroty
cd $(dirname $0)

# generate ca cert
openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt -subj /CN=ngrokex

# prepare ca directory
mkdir -p ca/certs && touch ca/index.txt && echo 01 > ca/serial

# generate & sign server cert
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -config server.cnf
openssl ca -batch -in server.csr -out server.crt -cert ca.crt -keyfile ca.key -days 1826 -config openssl.cnf

# generate & sign client cert
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -config client.cnf
openssl ca -batch -in client.csr -out client.crt -cert ca.crt -keyfile ca.key -days 1826 -config openssl.cnf

# copy cert files
cp ca.crt server.crt server.key ../assets/server/tls/
cp ca.crt client.crt client.key ../assets/client/tls/

# clean up
rm -r ca
rm *.csr

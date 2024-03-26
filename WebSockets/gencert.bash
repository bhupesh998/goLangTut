#!/bin/bash

echo "Creating Server.key"
openssl genrsa -out server.key 2048
oepnssl ecparam -genkey -name secp384r1 -out server.key
echo "Creating server.cert"
openssl req -new -x509 -sha256 -key server.key -out server.crt -batch -days 365
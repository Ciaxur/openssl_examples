#!/usr/bin/env bash

# Ensure working directory is in the right place.
CURRENT_DIR=$(realpath "$(dirname "$0")")
cd $CURRENT_DIR || { echo "Failed to change directory" && exit 1; }

# Generate private rsa key.
echo "Generating RSA key..."
openssl genrsa -out server.key 4096 || { echo "Failed to generate RSA key" && exit 1; }

# Generate CSR.
echo "Generating CSR..."
openssl req \
  -config ./server_v3conf.conf \
  -new \
  -key server.key \
  -out server.csr || { echo "Failed to generate CSR" && exit 1; }

# Generate a signed certificate.
echo "Generating signed certificate..."
openssl x509 \
  -req \
  -days 365 \
  -in ./server.csr \
  -CA ../nginx_server/ca/server_ca.crt \
  -CAkey ../nginx_server/ca/server_ca.key \
  -CAcreateserial \
  -extensions user_crt \
  -extfile ./server_v3conf.conf \
  -sha256 \
  -out server.crt || { echo "Failed to generate signed cert." && exit 1; }


#!/usr/bin/env bash

# Ensure working directory is in the right place.
CURRENT_DIR=$(realpath "$(dirname "$0")")
cd $CURRENT_DIR || { echo "Failed to change directory" && exit 1; }

# Generate private rsa key.
echo "Generating RSA key..."
openssl genrsa -out client.key 4096 || { echo "Failed to generate RSA key" && exit 1; }

# Generate CSR.
echo "Generating CSR..."
openssl req \
  -config ./client_v3conf.conf \
  -new \
  -key client.key \
  -out client.csr || { echo "Failed to generate CSR" && exit 1; }

# Generate a signed certificate.
echo "Generating signed certificate..."
openssl x509 \
  -req \
  -days 365 \
  -in ./client.csr \
  -CA ../nginx_client/ca/client_ca.crt \
  -CAkey ../nginx_client/ca/client_ca.key \
  -extensions user_crt \
  -extfile ./client_v3conf.conf \
  -CAcreateserial \
  -sha256 \
  -out client.crt || { echo "Failed to generate signed cert." && exit 1; }


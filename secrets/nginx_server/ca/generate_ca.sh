#!/usr/bin/env bash

# Ensure working directory is in the right place.
CURRENT_DIR=$(realpath "$(dirname "$0")")
cd $CURRENT_DIR || { echo "Failed to change directory" && exit 1; }

# Generate private rsa key.
openssl genrsa -out server_ca.key 4096 || { echo "Failed to generate RSA key" && exit 1; }

# Generate self-signed certificate.
echo "Generating self-signed CA..."
openssl req \
  -new \
  -x509 \
  -days 3650 \
  -config ./server_ca_v3conf.conf \
  -key ./server_ca.key \
  -out ./server_ca.crt || { echo "Failed to generate certificate" && exit 1; }


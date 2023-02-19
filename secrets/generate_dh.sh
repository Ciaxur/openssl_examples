#!/usr/bin/env bash

set -x

openssl dhparam -out dh.pem 4096 || { echo "Failed to generate dhparam" && exit 1; }
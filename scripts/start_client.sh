#!/usr/bin/env bash

set -x

go run ./client \
  --port 9000 \
  --cert ./secrets/go_client/client.crt \
  --key ./secrets/go_client/client.key \
  --ca ./secrets/nginx_server/ca/server_ca.crt

#!/usr/bin/env bash

set -x

go run ./server \
  --port 9001 \
  --cert ./secrets/go_server/server.crt \
  --key ./secrets/go_server/server.key \
  --ca ./secrets/nginx_client/ca/client_ca.crt

#!/usr/bin/env bash

curl --header "Content-Type: application/json" \
  --request PUT \
  --data "$(cat "$(dirname "$0")/tweets/index.json")" \
  http://localhost:9200/tweets
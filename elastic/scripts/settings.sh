#!/usr/bin/env bash

curl --request POST http://localhost:9200/tweets/_close

curl --header "Content-Type: application/json" \
  --request PUT \
  --data "$(cat "$(dirname "$0")/tweets/settings.json")" \
  http://localhost:9200/tweets/_settings

curl --request POST http://localhost:9200/tweets/_open
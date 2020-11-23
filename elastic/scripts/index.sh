#!/usr/bin/env bash

curl --header "Content-Type: application/json" \
  --request PUT \
  --data "$(cat tweets/mapping.json)" \
  http://localhost:9200/tweets/_mapping
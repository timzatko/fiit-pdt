#!/usr/bin/env bash

curl --header "Content-Type: application/json" \
  --request PUT \
  --data "$(cat tweets/indec.json)" \
  http://localhost:9200/tweets
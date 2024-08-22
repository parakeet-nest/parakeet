#!/bin/bash
set -o allexport; source .env; set +o allexport

: <<'COMMENT'
Add these variables to the .env file:
ELASTIC_USERNAME=elastic
ELASTIC_PASSWORD=<password>
ELASTIC_ADDRESS=https://localhost:9200
ELASTIC_CERT_PATH=../certs/ca/ca.crt
COMMENT

go run main.go


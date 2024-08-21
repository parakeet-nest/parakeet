#!/bin/bash
set -o allexport; source .env; set +o allexport

: <<'COMMENT'
create a .env file
add the ELASTIC_PASSWORD variable
COMMENT

curl --cacert http_ca.crt -u elastic:$ELASTIC_PASSWORD https://localhost:9200

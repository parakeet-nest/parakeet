#!/bin/bash
#set -o allexport; source release.env; set +o allexport

: <<'COMMENT'

COMMENT

docker cp es01:/usr/share/elasticsearch/config/certs/http_ca.crt .

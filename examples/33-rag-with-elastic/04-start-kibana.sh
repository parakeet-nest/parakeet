#!/bin/bash
#set -o allexport; source release.env; set +o allexport

: <<'COMMENT'

COMMENT

docker run --name kib01 --net elastic -p 5601:5601 docker.elastic.co/kibana/kibana:8.15.0
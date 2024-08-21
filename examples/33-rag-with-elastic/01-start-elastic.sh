#!/bin/bash
#set -o allexport; source release.env; set +o allexport

: <<'COMMENT'

COMMENT

docker network create elastic
docker run --name es01 --net elastic -p 9200:9200 -it -m 1GB docker.elastic.co/elasticsearch/elasticsearch:8.15.0


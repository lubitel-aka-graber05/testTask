#!/usr/bin/env bash

source "scripts/config/config.sh"

docker network create -d bridge --subnet=172.20.10.0/25 $NETWORKNAME || \
    echo "Network allready created"

docker container run --name db --network $NETWORKNAME -d -p :5432 -e POSTGRES_PASSWORD=password $POSTGRES || \
    docker container start db

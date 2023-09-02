#!/usr/bin/env bash

source "scripts/config/config.sh"

source "scripts/start_db.sh"

until docker container run --rm -d --name server --network $NETWORKNAME -p 51002:55005 $APPNAME; do \
    source "scripts/server_build.sh"
done







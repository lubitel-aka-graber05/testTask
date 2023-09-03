#!/usr/bin/env bash

POSTGRES="postgres:15.4-alpine3.18"

docker container run --rm -d --name testdb -e POSTGRES_PASSWORD=password -p 55001:5432 $POSTGRES
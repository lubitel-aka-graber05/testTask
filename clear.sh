#!/usr/bin/env bash

NETWORKNAME='testtasknetwork'
APPNAME='testtask'

docker container stop db
docker container stop server
docker network rm $NETWORKNAME
docker image rm $APPNAME
yes | docker image prune
#!/usr/bin/env bash

POSTGRES='postgres:15.4-alpine3.18'
ALPINE='alpine:latest'
NETWORKNAME='testtasknetwork'
APPNAME='testtask'

# docker pull $POSTGRES || { echo "Can't get postgres" ; exit ; }
# [[ docker pull $ALPINE ]] || { echo "Can't get alpine" ; exit ; } 

docker network create -d bridge --subnet=172.20.10.0/25 $NETWORKNAME || \
{ echo "Can't create docker network" ; exit ; }



docker container run --name db --network $NETWORKNAME --rm -d -p :5432 -e POSTGRES_PASSWORD=password $POSTGRES || \
{ echo "Can't run postgres container" ; exit ; }

docker build -t $APPNAME . || \
{ echo "Can't build image ($APPNAME)" ; exit ; }

docker container run --rm -d --name server --network $NETWORKNAME -p 51002:55005 $APPNAME || { echo "Can't run container" ; exit ; }




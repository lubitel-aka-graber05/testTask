#!/usr/bin/env bash

source "scripts/config/config.sh"

docker container stop $APPNAME
sleep 2s
source "start_server.sh" 
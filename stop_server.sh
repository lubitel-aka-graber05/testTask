#!/usr/bin/env bash


docker container stop server
sleep 1s
source "scripts/stop_db.sh"
#!/usr/bin/env bash

curl -X POST 0.0.0.0:51002/adduser -H 'Content-Type: application/json' -d '{"username":"testuser","password":"supersecret"}' -H 'Accept: application/json' || \
{ echo "Can't add auth info" ; exit ; } && \
echo "Add auth info completed"

curl -X POST 0.0.0.0:51002/addnote  -H 'Content-Type: application/json' -u "testuser:supersecret" -d '{"username":"testuser","body":"Hi!"}' -H 'Accept: application/json' || \
{ echo "Can't add note with basic auth"  ; exit ; } && \
echo "Add note completed"

curl -X POST 0.0.0.0:51002/outbyname  -H 'Content-Type: application/json' -u "testuser:supersecret" -d '{"username":"testuser"}' -H 'Accept: application/json' || \
{ echo "Can't output note with basic auth"  ; exit ; } && \
echo "Output note completed"


echo "Test container completed"

source "scripts/clear.sh"
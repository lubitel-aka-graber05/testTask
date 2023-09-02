#!/usr/bin/env bash

docker exec -it db psql -U postgres -d postgres -c "DROP TABLE note"

yes | docker image prune
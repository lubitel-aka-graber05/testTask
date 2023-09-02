#!/usr/bin/env bash

source "config/config.sh"


docker build -t $APPNAME . || \
{ echo "Can't build image ($APPNAME)" ; exit ; }


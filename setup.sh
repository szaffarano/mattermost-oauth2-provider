#!/bin/bash

DOCKER_UID=2000
DOCKER_GID=2000
BASE="./volumes/app/mattermost"

for d in data logs config plugins
do
  DIR="${BASE}/${d}"
  [ -d $DIR ] || mkdir -p $DIR
done

sudo chown -R $DOCKER_UID:$DOCKER_GID $BASE

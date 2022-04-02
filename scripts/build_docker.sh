#!/bin/sh

docker build -t iso-plugin -f ./docker/isoplugin/Dockerfile .
docker build -t iso-server -f ./docker/isoserver/Dockerfile .
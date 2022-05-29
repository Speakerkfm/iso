#!/bin/sh

docker build -t speakerkfm/iso-plugin -f ./docker/isoplugin/Dockerfile .
docker build -t speakerkfm/iso-server -f ./docker/isoserver/Dockerfile .
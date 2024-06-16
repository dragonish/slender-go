#!/bin/bash

docker image save giterhub/slender:latest | gzip > slender.tar.gz

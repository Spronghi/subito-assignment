#!/bin/sh
docker build -t subito-assignment .
docker run -p 8080:8080 subito-assignment


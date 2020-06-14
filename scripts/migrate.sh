#!/usr/bin/env bash

source .env
docker run -it -v /home/danko/apps/ecommerce/server/migrations:/migrations --network server_net migrate/migrate -path=/migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@database:5432/${POSTGRES_DB}?sslmode=disable" $@

#!/usr/bin/env bash

source .env

docker run -it -v /home/danko/apps/ecommerce/server/migrations:/migrations --network server_ecommerce-net migrate/migrate -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@database:5432/${POSTGRES_DB}?sslmode=disable" create -ext sql -dir /migrations -seq -digits 3 ${@
var/lib/postgresql/data
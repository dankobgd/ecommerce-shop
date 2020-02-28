#!/usr/bin/env bash

docker run -it -v /home/danko/apps/ecommerce/server/migrations:/migrations --network server_net migrate/migrate -database "postgres://test:test@database:5432/ecommerce?sslmode=disable" create -ext sql -dir /migrations -seq -digits 3 $@
var/lib/postgresql/data
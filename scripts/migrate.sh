#!/usr/bin/env bash

docker run -it -v /home/danko/apps/ecommerce/server/migrations:/migrations --network server_net migrate/migrate -path=/migrations -database "postgres://test:test@database:5432/ecommerce?sslmode=disable" $@

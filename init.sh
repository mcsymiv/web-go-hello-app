#!/usr/bin/env bash
#
# Starts application

if [ -f ./web ]; then
		echo "removing previous build web"
		rm ./web
fi

go build -o ./bin ./cmd/web

./bin/web \
		-prod=false \
		-cache=false \
		-addr=:8099 \
		-dbhost=localhost \
		-dbuser=postgres \
		-dbname=db \
		-dbpass="$1" \
		-dbport=5432

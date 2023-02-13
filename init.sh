#!/usr/bin/env bash
#
# Starts application

if [ -f ./web ]; then
		echo "removing previous build web"
		rm ./web
fi

go build ./cmd/web

./web \
		-prod=false \
		-cache=false \
		-addr=:8080 \
		-dbhost=localhost \
		-dbuser=postgres \
		-dbname=db \
		-dbpass="$1" \
		-dbport=5432

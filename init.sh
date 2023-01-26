#!/usr/bin/env bash
#
# Starts application
#

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

#!/bin/sh
set -e

# delete old db and create new one based on schema.sql
rm -f /app/data/mydb.sqlite
sqlite3 /app/data/mydb.sqlite < /app/data/schema.sql

# start the server
./server

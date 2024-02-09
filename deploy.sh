#! /bin/bash
echo "Running server config script..."

loc=$(pwd)

time=$(date)

mkdir -p "$loc/logs/$time"

sed -i "s/\"server_port\":.*/\"server_port\": \"${SERVER_INNER}\"/g" config.json

buildLog="$loc/logs/$time/build.log"
runLog="$loc/logs/$time/run.log"

echo "Building..."
go build >> "$buildLog"
echo "Done."

echo "Starting server..."
go run main.go >> "$runLog"
echo "Server stopped."
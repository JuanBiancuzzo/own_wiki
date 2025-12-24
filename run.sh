#!/usr/bin/env bash

echo "Building..."
go build -o ./bin/userDefineData.exe ./plugin
go build -o ./bin/app.exe ./exe

echo "Running app:"
./bin/app.exe
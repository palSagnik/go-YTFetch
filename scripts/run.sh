#!/bin/zsh

rm test
source .env
go mod tidy

echo "Building the backend"
go build -o test

echo "Running the project"
./test

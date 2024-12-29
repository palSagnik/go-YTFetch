#!/bin/zsh

rm test
source .env
go mod tidy

cd frontend
npm run dev&

cd ../
source .env
echo "Building the backend"
go build -o test

echo "Running the project"
./test

#! /usr/bin/env sh

set -e 

ROOT=$(pwd)

go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow@latest

echo "------------------- Checking API service"
cd $ROOT/services/api

echo "go test:"
go test --count=1 ./...
echo "go vet:"
go vet ./...
echo "staticcheck:"
staticcheck
echo "shadow:"
shadow ./...

cd $ROOT

echo "------------------- Checking common library"
cd $ROOT/services/common

echo "go test:"
go test --count=1 ./...
echo "go vet:"
go vet ./...
echo "shadow:"
shadow ./...

cd $ROOT

echo "------------------- Checking IMDB service"
cd $ROOT/services/imdb

echo "go test:"
go test --count=1 ./...
echo "go vet:"
go vet ./...
echo "staticcheck:"
staticcheck
echo "shadow:"
shadow ./...

cd $ROOT

echo "------------------- Checking Scraper service"
cd $ROOT/services/scraper

echo "go test:"
go test --count=1 ./...
echo "go vet:"
go vet ./...
echo "staticcheck:"
staticcheck
echo "shadow:"
shadow ./...


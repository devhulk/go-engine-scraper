#!/bin/bash

env GOOS=darwin GOARCH=amd64 go build .
mv enginescraper bin/mac-enginescraper

env GOOS=linux GOARCH=amd64 go build .
mv enginescraper bin/linux-amd64-enginescraper

env GOOS=windows GOARCH=amd64 go build .
mv enginescraper.exe bin/windows-a64-enginescraper.exe

env GOOS=windows GOARCH=386 go build .
mv enginescraper.exe bin/windows-386-enginescraper.exe

#!/bin/bash

name=ubx

# Build binary file for linux and mac
go build -o $name *.go

# Build binary file for windows
GOOS=windows GOARCH=386 go build -o $name.exe *.go

# Remove dist directory and re-create it
rm -rf ../dist
mkdir ../dist

# Move binary file to dist directory
mv $name ../dist
mv $name.exe ../dist

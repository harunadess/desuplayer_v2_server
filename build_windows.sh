#!/bin/bash
 env GOOS=windows GOARCH=amd64 go build -o server.exe main.go
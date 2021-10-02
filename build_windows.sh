#!/bin/bash
 env GOOS=windows GOARCH=amd64 go build
 mv "desuplayer_v2.exe" "desuplayer_v2 server.exe"
# /bin/bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
mv main ./../../build/mig-linux

go build main.go
mv main ./../../build/mig-mac
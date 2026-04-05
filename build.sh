GOOS=windows GOARCH=amd64 go build -o bin/windows/kee.exe kee.go
GOOS=linux GOARCH=amd64 go build -o bin/linux/kee kee.go
GOOS=darwin GOARCH=amd64 go build -o bin/macos/kee kee.go
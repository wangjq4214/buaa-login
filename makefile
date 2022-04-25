VERSION=0.2.0
PREFIX=login

build:
	go build -ldflags "-X main.Version=${VERSION}" -o ./bin/ ./cmd
 
run:
	go run ./cmd

test:
	go test -v ./...
 
compile:
	# Linux
	GOOS=linux GOARCH=386 go build -ldflags "-w -s -X main.Version=${VERSION}" -o bin/${PREFIX}-linux-386 ./cmd
	# Windows
	GOOS=windows GOARCH=386 go build -ldflags "-w -s -X main.Version=${VERSION}" -o bin/${PREFIX}-windows-386.exe ./cmd
	# 64-Bit
	# MacOS
	GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s -X main.Version=${VERSION}" -o bin/${PREFIX}-darwin-amd64 ./cmd
	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.Version=${VERSION}" -o bin/${PREFIX}-linux-amd64 ./cmd
	# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags "-w -s -X main.Version=${VERSION}" -o bin/${PREFIX}-windows-amd64.exe ./cmd
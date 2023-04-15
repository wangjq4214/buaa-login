VERSION=0.3.0
PREFIX=login

.PHONY: build
build:
	go build -ldflags "-X main.Version=${VERSION}" -o ./bin/ ./cmd

.PHONY: clean
clean:
	rm -r bin

.PHONY: test
test:
	go test -v ./...
 
.PHONY: compile
compile:
	# Windows
	GOOS=windows GOARCH=386 go build -ldflags "-w -s -X main.Version=${VERSION}" -o bin/${PREFIX}-windows-386.exe ./cmd
	GOOS=windows GOARCH=amd64 go build -ldflags "-w -s -X main.Version=${VERSION}" -o bin/${PREFIX}-windows-amd64.exe ./cmd
	# MacOS
	GOOS=darwin GOARCH=amd64 go build -ldflags "-w -s -X main.Version=${VERSION}" -o bin/${PREFIX}-darwin-amd64 ./cmd
	# Linux
	GOOS=linux GOARCH=386 go build -ldflags "-w -s -X main.Version=${VERSION}" -o bin/${PREFIX}-linux-386 ./cmd
	GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.Version=${VERSION}" -o bin/${PREFIX}-linux-amd64 ./cmd
	GOOS=linux GOARCH=arm64 go build -ldflags "-w -s -X main.Version=${VERSION}" -o bin/${PREFIX}-linux-arm64 ./cmd

clean:
	rm -Rfv bin
	mkdir bin

build: clean
	go build -o bin/restic-exporter src/*.go

build-all: clean
	GOOS="linux"   GOARCH="amd64"     go build -o bin/restic-exporter__linux-amd64 src/*.go
	GOOS="linux"   GOARCH="arm64"     go build -o bin/restic-exporter__linux-arm64   src/*.go
	GOOS="freebsd" GOARCH="amd64"     go build -o bin/restic-exporter__freebsd-amd64 src/*.go
	GOOS="freebsd" GOARCH="arm64"     go build -o bin/restic-exporter__freebsd-arm64 src/*.go

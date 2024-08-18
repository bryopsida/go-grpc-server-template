clean:
	rm -rf bin/*
build:
	go build -o bin/service main.go

image:
	docker build -t ghcr.io/bryopsida/go-grpc-server-template:local .

test:
	go test -v ./...
	
lint:
	golangci-lint run
	go install golang.org/x/lint/golint@latest
	golint ./...
	go vet ./...
	
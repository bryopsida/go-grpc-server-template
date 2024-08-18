clean:
	rm -rf bin/*

generate-grpc-code:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/v1/*.proto

build:
	go build -o bin/service main.go

image:
	docker build -t ghcr.io/bryopsida/go-grpc-server-template:local .

test:
	go test -v ./...
	
lint:
	go install golang.org/x/lint/golint@latest
	golint ./...
	go vet ./...

build:
	CGO_ENABLED=0 go build ./cmd/...

test:
	go test ./...

docker: build
	docker build . -t tap2junit:latest


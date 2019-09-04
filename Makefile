VERSION := "v0.0.7"

build:
	CGO_ENABLED=0 go build ./cmd/...

test:
	go test ./...

docker: build
	docker build . -t tap2junit:latest

tag: docker
	docker tag tap2junit:latest filipfilmar/tap2junit:${VERSION}

# This bit requires a valid docker login
push: tag
	docker push filipfilmar/tap2junit:${VERSION}


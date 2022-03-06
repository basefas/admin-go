.PHONY: run build

# $(shell git rev-parse --short HEAD)
VERSION := 0.0.3

run:
	go run ./cmd/app/main.go

build:
	go build -o ./build/bin/admin-go -v ./cmd/app/main.go

docker-build:
	docker build \
		-f deploy/docker/Dockerfile \
		-t admin-go:$(VERSION)-amd64 \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

docker-build-chn:
	docker build \
		-f deploy/docker/chn.Dockerfile \
		-t admin-go:$(VERSION)-amd64-chn \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.
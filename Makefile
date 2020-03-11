test:
	go test ./...

fmt:
	gofmt -w .

download:
	go mod download

.PHONY: download fmt test

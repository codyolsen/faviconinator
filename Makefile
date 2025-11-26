BINARY := faviconinator
CMD := ./cmd/faviconinator

.PHONY: build install test fmt vet clean

build:
	go build -o bin/$(BINARY) $(CMD)

install:
	go install $(CMD)

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

clean:
	rm -rf bin

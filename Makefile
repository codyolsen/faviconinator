BINARY := faviconinator
CMD := ./cmd/faviconinator
TARGETS := mac/amd64 mac/arm64 linux/amd64 linux/arm64 windows/amd64

.PHONY: build build-release install test fmt vet clean

build:
	go build -o bin/$(BINARY) $(CMD)

build-release:
	@mkdir -p bin
	@for target in $(TARGETS); do \
		OS_NAME=$${target%/*}; GOARCH=$${target#*/}; \
		case "$$OS_NAME" in \
			mac) GOOS=darwin ;; \
			linux) GOOS=linux ;; \
			windows) GOOS=windows ;; \
			*) echo "Unknown OS $$OS_NAME"; exit 1 ;; \
		esac; \
		ext=""; [ "$${GOOS}" = "windows" ] && ext=".exe"; \
		out=bin/$(BINARY)-$${OS_NAME}-x$${GOARCH}$${ext}; \
		echo "GOOS=$$GOOS GOARCH=$$GOARCH -> $$out"; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build -o $$out $(CMD); \
	done

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

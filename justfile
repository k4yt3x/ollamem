build:
	go build -ldflags="-s -w" -trimpath -o bin/ollamem ./cmd/ollamem

debug:
	go build -o bin/ollamem ./cmd/ollamem

test:
    go test -v ./...

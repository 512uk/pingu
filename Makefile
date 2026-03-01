.PHONY: build run test lint clean

build:
	go build -o bin/pingu ./cmd/pingu

run:
	go run ./cmd/pingu

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -rf bin/

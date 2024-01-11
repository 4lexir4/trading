build:
	go build -o bin/trading

run: build
	./bin/trading

test:
	go test -v ./...

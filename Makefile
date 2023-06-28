build:
	go build -o bin/blockchain-layer
run: build
	./bin/blockchain-layer
test:
	go test ./...
BINARY_NAME=fjroland
clean:
	rm -f bin/$(BINARY_NAME)
build:
	go build -o bin/$(BINARY_NAME) cmd/fjroland/main.go

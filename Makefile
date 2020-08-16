BINARY_NAME=fjroland
clean:
	rm -f bin/$(BINARY_NAME)
build:
	go build -o bin/$(BINARY_NAME) cmd/fjroland/main.go

demo: build
	bin/fjroland assets/patterns/dont-say-nuthin.json

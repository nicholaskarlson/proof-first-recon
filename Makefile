.PHONY: test demo build clean

test:
	go test ./...

demo:
	go run ./cmd/recon -mode demo -out ./out

build:
	mkdir -p bin
	go build -o bin/recon ./cmd/recon

clean:
	rm -rf ./out ./bin

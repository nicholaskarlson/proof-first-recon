.PHONY: test demo build clean

test:
	go test -count=1 ./...

demo:
	go run ./cmd/recon demo --out ./out

build:
	mkdir -p bin
	go build -o bin/recon ./cmd/recon

clean:
	rm -rf ./out ./bin

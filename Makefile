.PHONY: test demo verify build clean

test:
	go test -count=1 ./...

demo:
	go run ./cmd/recon demo --out ./out

# Book-facing proof gate: tests + demo fixture verification
verify: test demo

build:
	mkdir -p bin
	go build -o bin/recon ./cmd/recon

clean:
	rm -rf ./out ./bin

BINARY := nuxtblog

.PHONY: build install lint snapshot clean

build:
	go build -o $(BINARY) ./cmd/nuxtblog

install:
	go install ./cmd/nuxtblog

lint:
	go vet ./...

snapshot:
	goreleaser release --snapshot --clean

clean:
	rm -f $(BINARY)

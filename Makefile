BINARY := nuxtblog

.PHONY: build install lint snapshot clean

build:
	go build -o $(BINARY) .

install:
	go install .

lint:
	go vet ./...

snapshot:
	goreleaser release --snapshot --clean

clean:
	rm -f $(BINARY)

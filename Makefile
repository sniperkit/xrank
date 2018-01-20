default: build

build:
	go build

test:
	go test -tags mocks ./...

clean:
	go clean

install:
	go install

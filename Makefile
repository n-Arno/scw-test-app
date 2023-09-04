all: build

build:
	GOOS=linux GOARCH=amd64 go build -o scw-test-app-linux-amd64

test:
	go build && ./scw-test-app

clean:
	go clean

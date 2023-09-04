all: build

build:
	go mod tidy && go build


clean:
	go clean

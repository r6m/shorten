IMAGE=r6m/shorten

all: clean build

clean:
	rm -f bin/*

build:
	GOOS=linux GOARCH=amd64 go build -o bin/shorten .

docker:
	docker build -t ${IMAGE}:v0.1.0
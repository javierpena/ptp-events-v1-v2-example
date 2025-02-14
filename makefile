all: deps build image push

clean:
	rm -f bin/*
	rm -rf vendor
	rm -f go.sum

build: cmd/main.go
	go build -o ./bin/simpleapp cmd/main.go

image: bin/simpleapp
	podman build -t quay.io/javierpena/consumerapp:0.1 .

push:
	podman push quay.io/javierpena/consumerapp:0.1

deps:
	go mod tidy && \
	go mod vendor

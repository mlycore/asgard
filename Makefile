shell := bash

clean:
	rm ./bin/asgard

.PHONY: build
build:
	go build -o ./bin/asgard ./pkg

run:
	./bin/asgard

registry:
	docker build -f ./build/Dockerfile . -t mworks92/asgard:latest
	docker push mworks92/asgard:latest

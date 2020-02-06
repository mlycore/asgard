shell := bash

clean:
	rm ./bin/asgard

build:
	go build -o ./bin/asgard ./pkg

run:
	./bin/asgard

registry:
	docker build -f ./deployments/Dockerfile . -t mworks92/asgard:latest
	docker push mworks92/asgard:latest

shell := bash

build:
	go build -o ./out/asgard ./pkg

run:
	make build
	./out/asgard

registry:
	docker build -f ./deployments/Dockerfile . -t mworks92/asgard:latest
	docker push mworks92/asgard:latest

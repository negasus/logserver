TAG := $(shell git describe --tag)

.PHONY: build push

build:
	@echo Build $(TAG)
	docker build --build-arg version=$(TAG) -t negasus/logserver:$(TAG) -f Dockerfile .
push:
	@echo Push $(TAG)
	docker push negasus/logserver:$(TAG)
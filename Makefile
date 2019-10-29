TAG := $(shell git describe --tag)

.PHONY: build push

build:
	@echo Build $(TAG)
	docker build -t negasus/logserver:$(TAG) -f Dockerfile .
push:
	@echo Push $(TAG)
	docker push negasus/logserver:$(TAG)
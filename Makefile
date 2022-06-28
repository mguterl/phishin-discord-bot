PROJECT_NAME := "mguterl/phishin-discord-bot"

.PHONY: docker-build publish publish-latest publish-version fetch-version

release: docker-build publish

docker-build: fetch-version
	docker build --platform linux/amd64 -t $(PROJECT_NAME):latest -t $(PROJECT_NAME):$(VERSION) .

publish: publish-latest publish-version

publish-latest:
	docker push $(PROJECT_NAME):latest

publish-version: fetch-version
	docker push $(PROJECT_NAME):$(VERSION)

fetch-version:
	$(eval VERSION := $(shell git rev-parse --short HEAD))

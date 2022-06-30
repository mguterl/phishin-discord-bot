PROJECT_NAME := "mguterl/phishin-discord-bot"

.PHONY: docker-build publish publish-latest publish-version fetch-version

release: docker-build publish

docker-build:
	docker build --platform linux/amd64 -t $(PROJECT_NAME) .

publish:
	docker push $(PROJECT_NAME)

APP_NAME = log-generator

DOCKER_REPO = rea1shane
DOCKER_IMAGE_NAME = $(APP_NAME)

.PHONY: docker-run
docker-run: docker-build
	docker run --name $(APP_NAME) -d $(DOCKER_REPO)/$(DOCKER_IMAGE_NAME)

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_REPO)/$(DOCKER_IMAGE_NAME) .

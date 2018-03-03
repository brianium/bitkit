.PHONY: docker-login
docker-login:
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASSWORD)

.PHONY: server-image-build
server-image-build:
	docker build -t scaturr/memcool .

.PHONY: server-push
server-image-push:
	docker push scaturr/memcool

.PHONY: server-image
server-image: server-image-build server-image-push

.PHONY: deploy
deploy: docker-login server-image

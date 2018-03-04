.PHONY: docker-login
docker-login:
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASSWORD)

.PHONY: server-image-build
server-image-build:
	docker build -t scaturr/memcool server

.PHONY: server-push
server-image-push:
	docker push scaturr/memcool

.PHONY: server-image
server-image: server-image-build server-image-push

.PHONY: deploy
deploy: docker-login server-image

.PHONY: run
run:
	docker-compose up --build -d

.PHONY: stop
stop:
	docker-compose down

.PHONY: restart
restart: stop run

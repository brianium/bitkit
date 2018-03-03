.PHONY: server-image
server-image-build:
	docker build --iidfile server.img .

.PHONY: server-push
server-image-push:
	docker tag $(shell cat server.img) scaturr/memcool
	docker push scaturr/memcool
	rm server.img

.PHONY: server-image
server-image: server-image-build server-image-push

.PHONY: deploy
deploy: server-image

# certs command
gencerts = openssl req \
    -newkey rsa:2048 \
    -x509 \
    -nodes \
    -keyout server/key.pem \
    -new \
    -out server/cert.pem \
    -subj /CN=localhost \
    -sha256 \
    -days 3650

# setup the stage
ifeq ("$(CIRCLE_BRANCH)", "master")
	STAGE = production
else
	STAGE = staging
endif

ifeq ("$(ENV)", "development")
	STAGE = development
endif

# Set up commands based on stage
ifeq ("$(STAGE)","production")
	prebuild = echo "Building production image"
	build = docker build -t scaturr/bitkit server
	push = docker push scaturr/bitkit
	pguri = $(POSTGRES_URI)
else
	prebuild = $(gencerts)
	build = docker build -t scaturr/bitkit:staging server
	push = docker push scaturr/bitkit:staging
	pguri = $(POSTGRES_STAGING_URI)
endif

ifeq ("$(STAGE)","development")
	pguri = $(subst db,localhost,$(POSTGRES_URI))
endif

.PHONY: docker-login
docker-login:
	docker login -u $(DOCKER_USER) -p $(DOCKER_PASSWORD)

.PHONY: server-image-build
server-image-build:
	$(prebuild)
	$(build)

.PHONY: server-image-push
server-image-push:
	$(push)

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

.PHONY: migrate
migrate:
	pgmigrate -c $(pguri) -d db -t latest migrate

.PHONY: certs
certs:
	$(gencerts)

.PHONY: client
client:
	cd client/bitkit && lein clean && lein cljsbuild once min

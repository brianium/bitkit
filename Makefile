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
ifeq ("$(STAGE)","staging")
	pguri = $(POSTGRES_URI_STAGING)
else
	pguri = $(POSTGRES_URI)
endif

.PHONY: startdb
startdb:
	docker-compose up --build -d

.PHONY: stopdb
stopdb:
	docker-compose down

.PHONY: restartdb
restartdb: startdb stopdb

.PHONY: migrate
migrate:
	pgmigrate -c $(pguri) -d db -t latest migrate

.PHONY: certs
certs:
	$(gencerts)

.PHONY: client
client:
	cd client/bitkit && lein clean && lein cljsbuild once min

build:
	cd server && GOOS=linux go build -o main
	cd server && zip deployment.zip main
	mv server/deployment.zip .

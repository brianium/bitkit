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
	pgmigrate -c $(POSTGRES_URI) -d db -t latest migrate

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

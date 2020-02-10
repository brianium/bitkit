# bitkit

Formerly hosted at https://bitkit.live. Now defunct, but left here because it was a fun and cool project.

Highlights:

* A web server written in Go for use via AWS Lambda and API Gateway
* An app written in python for pushing bitcoin mempool data to the Go API
* A ClojureScript frontend for viewing transactions as they move through the mempool

## Requirements

* [Docker](https://www.docker.com/) (optional for dev)
* [Docker Compose](https://docs.docker.com/compose/) (optional for dev)
* [Make](https://www.gnu.org/software/make/)
* [direnv](https://direnv.net/)
* [python3](https://www.python.org/)
* [PGMigrate](https://github.com/yandex/pgmigrate)
* [Leiningen](https://leiningen.org/)
* Java 8 (For Leiningen)

## Setup

Set environment variables

```
$ direnv allow
```

To get a running db (with docker):

```
$ make startdb
```

You will need to run `direnv allow` any time the .envrc file changes

Make sure you change `.envrc.sample` to `.envrc` to fit your config needs

Run database migrations

```
$ make migrate
```

Migrations should be run after `make restart` as well. TODO - add `migrate` task to `restart`. (tricky because you must wait for postgres to be up)

You may want to setup `$GOPATH` correctly on your machine. This is optional, but may be
desired for local tooling. The docker container for the go app has a working dir of `/go/src/app`.
Creating a symlink from the `server` directory to your local `$GOPATH` should do the trick:

```
$ ln -s /Users/username/projects/bitkit/server /Users/username/go/src/server
```

Where `username` is your own user name. The path examples above are conventional for mac systems - so adjusting
for another system may be needed - i.e `/home/user/...` on a Linux system.

### SSL

A self signed cert is used for local development - and routes will only be accessible over `https`. To generate
local certs, run:

```
$ make certs
```

Fill out the asked questions. You need to run `make restart` if the docker containers are already up and running. You may
have to add an exception in your browser to make things work.

## Running

To run the API use `go run` within the server directory:

```
$ cd server && go run
```

Because Go programs are compiled, changes will not be reflected immediately. 

The go app exposes the web application on port 8080. You can visit the application
locally at `http://127.0.0.1:8080/`

## Client

The client application is written in ClojureScript using the [re-frame](https://github.com/Day8/re-frame) framework.

To start an interactive development environment, from within the `client/bitkit` directory
run:

```
$ lein figwheel
```

Production builds are done using:

```
$ make client
```

From the project root.

## Data

There are some outstanding field naming issues.

`fee_rate` is what I would prefer to call `mining_fee_rate`. It's the effective fee rate of the transaction considering ancestors and descendants as most miners would. Unless a transaction has no ancestors or descendants, it may be different than fee / vsize for the given transaction. Child Pays For Parent (CPFP) necessitates inclusion of parent transactions with or before child transactions.

`weight` is really `vsize` for the given transaction only (no ancestors or descendants).

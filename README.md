# memcool

## Requirements

* [Docker](https://www.docker.com/)
* [Docker Compose](https://docs.docker.com/compose/)
* [Make](https://www.gnu.org/software/make/)
* [direnv](https://direnv.net/)

## Setup

Set environment variables

```
$ direnv allow
```

You will need to run `direnv allow` any time the .envrc file changes

Make sure you change `.envrc.sample` to fit your config needs

## Running

To run the application use `docker-compose`:

```
$ docker-compose up
```

Or to run in the background

```
$ docker-compose up -d
```

Because Go programs are compiled, changes will not be reflected immediately. To see changes take effect you
may need to run `docker-compose` with the `--build` switch:

```
$ docker-compose up --build
```

The docker container exposes the web application on port 8080. You can visit the application
locally at `http://localhost:8080/`

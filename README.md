# memcool

## Requirements

* [Docker](https://www.docker.com/)
* [Docker Compose](https://docs.docker.com/compose/)
* [Make](https://www.gnu.org/software/make/)
* [direnv](https://direnv.net/)
<<<<<<< HEAD
* [python3](https://www.python.org/)
=======
* [PGMigrate](https://github.com/yandex/pgmigrate)
>>>>>>> Adds migrate target to Makefile

## Setup

Set environment variables

```
$ direnv allow
```

You will need to run `direnv allow` any time the .envrc file changes

Make sure you change `.envrc.sample` to `.envrc` to fit your config needs

## Running

To run the application use `make run`:

```
$ make run
```

Because Go programs are compiled, changes will not be reflected immediately. To see changes take effect you
need to run `docker-compose` with the `--build` switch. The `run` target in the Makefile handles this for you.

The docker container exposes the web application on port 8080. You can visit the application
locally at `http://localhost:8080/`

To stop use `make stop`. Or to stop, rebuild, and start again:

```
$ make restart
```

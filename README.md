# memcool

## Requirements

* [Docker](https://www.docker.com/)
* [Docker Compose](https://docs.docker.com/compose/)

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

# FX Template

> Go application template based on [Uber Fx](https://uber-go.github.io/fx).

## Usage

You first need to:
```shell
mv .env.example .env
```

Then, this project provides a [Makefile](Makefile), offering the following commands:

```shell
make up                   # start the docker compose stack
make down                 # stop the docker compose stack
make logs                 # stream docker compose stack logs
make fresh                # refresh docker compose stack
make build name={my_name} # build an application image named my_mame
```

# FX Template

> Go application template based on [Uber Fx](https://uber-go.github.io/fx).

## Usage

You first need to copy and adapt the `.env` file:
```shell
cp .env.example .env
```

Then, this project provides a [Makefile](Makefile), offering the following commands:

```shell
# app
make up                   # start the docker compose stack
make down                 # stop the docker compose stack
make logs                 # stream the docker compose stack logs
make fresh                # refresh the docker compose stack

# tools
make build name={my_img}  # build an application image named my_img
```

# product-sv

Microservice written in Go provides REST API to operate with product subdomain.

## Run Prerequisites

### Migrations

There is no auto migration setup for this app.

It's recommended to install and use goose migration tool.

You can find how to install and use goose in [Makefile](./cmd/migrator/Makefile).

---

### Optional Run Prerequisites

Optional prerequisites depend on the following options how you run app using:

1. Go (1.24.6 or upper)
2. Docker
3. Docker Compose

### env file

1. Create copy of [.env.example](../.env.example) file.
2. Set values depends on your environment.
3. Move .env file to current directory.

---

## Run docker container

### Build image
`
docker build --no-cache -f ./docker/Dockerfile -t product-sv:latest ../
`

#### Run container

`
docker run --env-file ./.env product-sv ../
`

---

## TODO

Future improvements:
 
1. Service layer performs sequence of unreliable request to external systems. 
It's recommended to add distributed transaction or alternative approach.
2. Add tests.

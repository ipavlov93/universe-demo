# product-sv

Microservice written in Go provides REST API to operate with product subdomain.

## Run Prerequisites

### env file for docker-compose run

1. Create copy of [.env.example](../.env.example) file.
2. Set values depends on your environment.
3. Move .env file to ./compose directory.

### Migrations

It's recommended to install and use [goose](https://github.com/pressly/goose) migration tool.

You can find commands how to install and run database migrations in [Makefile](./cmd/migrator/Makefile).

### Optional Run Prerequisites

Optional prerequisites depend on the following options how you run this app using:

1. Go (1.24.6 or upper)
2. Docker
3. Docker Compose

---

### Test API

You can find http API call examples here [product.http](./cmd/api/product.http).

You can easily change environment to another in [http-client.env.json](./cmd/api/http-client.env.json).

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
 
1. Service and Facade layer performs sequence of unreliable request to external systems. 
It's recommended to add distributed transaction or alternative approach with retry strategy in background.
2. Add tests.

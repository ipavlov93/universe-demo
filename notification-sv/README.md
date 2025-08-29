# notification-sv

Microservice written in Go designed as three-stage pipeline using workers and channel.

---

## Run Prerequisites

There are several options how you can run app using:

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
docker build --no-cache -f ./docker/Dockerfile -t notification-sv:latest ../
`

#### Run container

`
docker run --env-file ./.env notification-sv ../
`

---

## Concept

App designed as three-stage pipeline using workers and channel.
Buffered channels are used to prevent immediate block on channel send operation.

---

## TODO

Future improvements:
 
1. Consumer and Processor workers should have possibility to process remaining in-memory messages as part of graceful shutdown.
2. Add tests.
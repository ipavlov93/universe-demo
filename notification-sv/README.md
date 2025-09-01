# notification-sv

Microservice written in Go designed as three-stage pipeline using workers and channel.

---

## Run Prerequisites

### env file for docker-compose run

1. Create copy of [.env.example](../.env.example) file.
2. Set values depends on your environment.
3. Move .env file to ./compose directory.

### Optional Run Prerequisites

Optional prerequisites depend on the following options how you run this app using:

1. Go (1.24.6 or upper)
2. Docker
3. Docker Compose

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

### MessageProcessor

Represents a component that writes incoming messages to given io.Writer concurrently.
Send processed message's receiptHandles to next pipeline's stage. 

### Consumer

Represents a component with two pipeline's stages:
1. Consuming messages from external message broker concurrently.
2. Acknowledging message receive to message broker concurrently.

---

## TODO

Future improvements:
 
1. Add tests.
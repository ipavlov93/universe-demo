# universe-demo

universe-demo is mono repository.

universe-demo consists of Go modules: 

- notification-sv
- product-sv
- product-eventbus-pkg
- universe-pkg

## Demo branch

The actual branch for demo (with latest changes) `branch-for-demo`: https://github.com/ipavlov93/universe-demo/tree/branch-for-demo

---

## Repository Go modules

### notification-sv

[notification-sv documentation](./notification-sv/README.md)

### product-sv

[product-sv documentation](./product-sv/README.md)

### universe-pkg

Module contains packages with shared types and utility functions by other modules.

---

## Run Prerequisites

There are several options how you can run this mono repository's apps using:

1. Go (1.24.6 or upper)
2. Docker
3. Docker Compose

### env file

1. Create copy of [.env.example](.env.example) file.
2. Set values depends on your environment.
3. Move .env file to ./compose directory.

---

## Run cluster using Docker Compose

## docker-compose

You can find more commands in [Makefile](Makefile).

#### build and run
`
docker-compose -f ./compose/docker-compose.yml build && docker-compose -f ./compose/docker-compose.yml up -d
`

#### run
`
docker-compose -f ./compose/docker-compose.yml up -d
`

---

## Development


#### Generate mocks using mockery

`mockery --name=YourService --dir=your_path --output=your_path/mocks`

### Migration tool

Prerequisites: [goose](https://github.com/pressly/goose).

Installation:
`
go install github.com/pressly/goose/v3/cmd/goose@latest
`

Look at [Makefile](.product-sv/cmd/migrator/Makefile) that contains ready-to-use commands for running database migrations.

### Formatting

`go fmt ./...`

### goimports

To group and sort import sections example:

`goimports --local universe-demo -l -w .`

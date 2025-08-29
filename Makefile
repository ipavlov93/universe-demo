# docker-compose commands for local development
# Prerequisites: Docker, docker-compose.

docker_compose_build_and_run:
	docker-compose -f ./compose/docker-compose.yml build && docker-compose -f ./compose/docker-compose.yml up -d

docker_compose_build:
	docker-compose -f ./compose/docker-compose.yml build

docker_compose_run:
	docker-compose -f ./compose/docker-compose.yml up -d

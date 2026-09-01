CFG_DIR=./internal/config/local.yaml
CMD=./cmd/UserServer/main.go

docker_build:
	docker compose down
	docker compose build --no-cache
	docker-compose up -d

docker_run:
	docker-compose up -d

run:
	go run $(CMD) "-config=$(CFG_DIR)"
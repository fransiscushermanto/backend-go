build:
	go build -o bin/backend-api ./cmd/api

build-docker:
	docker compose build

dev:
	docker compose -f docker-compose.yaml -f docker-compose.dev.yaml $(action) $(if $(filter true,$(background)),-d,)  $*

setup-pkg:
	sh ./scripts/setup_dev_env.sh

seed:
	docker compose -f "docker-compose.yaml" -f docker-compose.dev.yaml run \
	--rm \
	api go run cmd/seeder/main.go $*

migrate:
	docker compose -f "docker-compose.yaml" run \
	--build \
	--rm \
	--entrypoint sh \
	api /app/scripts/migrate.sh $(action) $(version) $*

generate-key:
	@echo "🔑 Generating JWT keys..."
	docker compose -f "docker-compose.yaml" run \
	--rm \
	api go run cmd/keygen/main.go
	@echo "📁 Keys generated in project directory"
	@echo "🔧 Setting environment variables..."
	@echo "Run the following commands to set environment variables:"
	@echo "\nexport PRIVATE_KEY=\"$$(cat private_key.pem)\""
	@echo "\nexport PUBLIC_KEY=\"$$(cat public_key.pem)\""

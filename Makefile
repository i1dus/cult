ifeq ($(POSTGRES_SETUP),)
	POSTGRES_SETUP := user=postgres password=password dbname=cult host=localhost port=5432 sslmode=disable
endif

DATABASE_URL=postgresql://postgres:password@localhost:5434/hotel_management?sslmode=disable
MIGRATION_FOLDER=$(CURDIR)/migrations

db-up:
	docker-compose up

db-down:
	docker-compose down

migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create rename_me sql

.migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" up

.migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP)" down-to 0

migration-up: .migration-up
migration-down: .migration-down
db-reset: .migration-down .migration-up

run:
	CONFIG_PATH=./internal/config/config_local.yaml go run cmd/main.go

test:
	go test ./...

generate:
	protoc -I api api/parking_lot.proto --go_out=./internal/gen/parking_lot --go_opt=paths=source_relative --go-grpc_out=./internal/gen/parking_lot --go-grpc_opt=paths=source_relative
	protoc -I api api/sso.proto --go_out=./internal/gen/sso --go_opt=paths=source_relative --go-grpc_out=./internal/gen/sso --go-grpc_opt=paths=source_relative
.PHONY: proto migrate-up migrate-down proto migate-create

migrate-create:
	@test -n "$(NAME)" || (echo "Usage: make migrate-create NAME=migration_name" && exit 1)
	migrate create -ext sql -dir migrations $(NAME)

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

proto:
	protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/log.proto

include .env
export

.PHONY: migrate-up migrate-down migrate-create migrate-force

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_ADDR)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_ADDR)" down 1

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) $(name)

migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_ADDR)" force $(version)

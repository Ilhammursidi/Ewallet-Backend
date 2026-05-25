include ./.env.migration

MIGRATION_PATH=db/migrations
SEEDER_PATH=db/seeder
DATABASE_URL=postgresql://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
migrate-create:
	@migrate create -ext sql -dir $(MIGRATION_PATH) -seq create_$(NAME)_table

migrate-up:
	@migrate -database $(DATABASE_URL) -path $(MIGRATION_PATH) up

migrate-down:
	@migrate -database $(DATABASE_URL) -path $(MIGRATION_PATH) down

migrate-force:
	@migrate -database $(DATABASE_URL) -path $(MIGRATION_PATH) force $(VERSION)

# migrate-seed:
# 	@migrate -database $(DATABASE_URL) -path $(SEEDER_PATH) up

# migrate-seed-up:
# 	@migrate create -ext sql -dir $(SEEDER_PATH) -seq $(NAME)

# migrate-seed-down:
# 	@migrate -database $(DATABASE_URL) -path $(SEEDER_PATH) down

# migrate-seed-force:
# 	@migrate -database $(DATABASE_URL) -path $(SEEDER_PATH) force $(VERSION)

migrate-print:
	@echo $(DATABASE_URL)


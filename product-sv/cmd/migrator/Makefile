include ../../../compose/.env

GOOSE_ENV=GOOSE_MIGRATION_DIR=../../internal/migrations/postgres GOOSE_DRIVER=${GOOSE_DRIVER} GOOSE_DBSTRING=${GOOSE_DBSTRING}

# Makefile has list of following migrate commands using goose cli tool:
# Prerequisites: Go 1.16+, goose.
goose_install:
	go install github.com/pressly/goose/v3/cmd/goose@latest

# always create new migration file and apply sequential ordering to it (if previous files exist)"
create_valid_migration: create_migration migrate_fix migrate_validate

MIGRATION_NAME=$(or $(MIGRATION), new_migration)
create_migration:
	@echo "Running goose create init sql"
	$(GOOSE_ENV) goose create $(MIGRATION_NAME) sql

migrate_version:
	$(GOOSE_ENV) goose version

migrate_status:
	$(GOOSE_ENV) goose status

migrate_fix:
	@echo "Running goose fix to apply sequential ordering to migrations"
	$(GOOSE_ENV) goose fix

migrate_up:
	@echo "Running goose up to migrate the DB to the most recent version available"
	$(GOOSE_ENV) goose up

migrate_up_by_one:
	@echo "Running goose up to migrate the DB up-by-one"
	$(GOOSE_ENV) goose up-by-one

migrate_down_by_one:
	@echo "Running goose down to rollback the DB down-by-one"
	$(GOOSE_ENV) goose down

migrate_validate:
	@echo "Running goose validate migration files (goose syntax only) without running them"
	$(GOOSE_ENV) goose validate

.PHONY: app-start app-stop gen lint migrate-new migrate-up migrate-down migrate-drop

# App
app-start:
	docker-compose up -d --build --remove-orphans --force-recreate
app-stop:
	docker-compose down

# Gen
gen:
	go generate ./...

lint:
	golangci-lint run -v ./...

# Migrations
migrate-new:
	goose -dir ./migrations create $(name) sql && goose -dir ./migrations fix
migrate-up:
	goose -dir ./migrations postgres ${pg_dsn} up
migrate-down:
	goose -dir ./migrations postgres ${pg_dsn} down 1
migrate-drop:
	goose -dir ./migrations postgres ${pg_dsn} reset


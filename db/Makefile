postgres:
	docker run -d --name simpletodo -p 5432:5432 -e POSTGRES_USER=pasan -e POSTGRES_PASSWORD=12345 postgres:16-alpine

createdb:
	docker exec -it simpletodo createdb --username=pasan --owner=pasan simpletodo

create_migrations_up_down:
	migrate create -ext sql -dir db/migrations -seq init_mg

dropdb:
	docker exec -it simpletodo dropdb --username=pasan simpletodo

migrateup:
	migrate -path db/migrations -database "postgresql://pasan:12345@localhost:5432/simpletodo?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://pasan:12345@localhost:5432/simpletodo?sslmode=disable" -verbose up" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://pasan:12345@localhost:5432/simpletodo?sslmode=disable" -verbose up" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://pasan:12345@localhost:5432/simpletodo?sslmode=disable" -verbose up" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq simpletodo

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

sqlc:
	sqlc generate

.PHONY: postgres createdb migrations dropdb migrateup migratedown sqlc
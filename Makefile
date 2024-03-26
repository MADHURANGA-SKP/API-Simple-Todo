DB_URL = postgresql://pasan:12345@localhost:5432/simpletodo?sslmode=disable

postgres:
	docker run -d --name simpletodo -p 5432:5432 -e POSTGRES_USER=pasan -e POSTGRES_PASSWORD=12345 postgres:16-alpine

createdb:
	docker exec -it simpletodo createdb --username=pasan --owner=pasan simpletodo

create_migration_up_down:
	migrate create -ext sql -dir db/migration -seq init_mg

dropdb:
	docker exec -it simpletodo dropdb --username=pasan simpletodo

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up


migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose up" -verbose down


new_migration: 
	migrate create -ext sql -dir db/migration -seq simpletodo

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simpletodo \
	--experimental_allow_proto3_optional \
	proto/*.proto
	statik -src=./doc/swagger -dest=./doc

rm_proto:
	rm -f pb/*.go

redis:
	docker run --name redis -p 6379:6379 -d redis:alpine3.19

server:
	go run main.go

.PHONY: postgres createdb create_migration_up_down dropdb migrateup migratedown db_docs db_schema sqlc test proto rm_proto server redis
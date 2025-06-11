DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratredown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
db_docs:
	dbdocs build doc/db.dbml
db_schema:
	dbml2sql -o db/schema.sql doc/db.dbml
sqlc:
	sqlc generate

test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go SimpleBank/db/sqlc Store

proto:
	rm -rf pb/*.go
	protoc --proto_path=proto \
		--go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
		proto/*.proto
evans:
	evans --host localhost --port 9090 -r repl
.PHONY: createdb dropdb postgres migratedown migrateup migratedown1 migrateup1 db_docs db_schema sqlc test server mock proto evans
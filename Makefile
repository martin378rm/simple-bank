migrateup: 
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown: 
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose down


postgres: 
	docker run --name postgres_test -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:15.5

createdb:
	docker exec -it postgres_test createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it postgres_test dropdb simple_bank

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc
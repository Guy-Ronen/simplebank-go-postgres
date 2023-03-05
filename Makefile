postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=guy.ronen -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=guy.ronen simple_bank

dropdb:
	docker exec -it postgres12 dropdb --username=guy.ronen simple_bank

rundb:
	docker exec -it postgres12 psql -U guy.ronen simple_bank

migreateup:
	migrate -path db/migration -database "postgresql://guy.ronen:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://guy.ronen:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migreateup migratedown sqlc
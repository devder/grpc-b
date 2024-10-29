DB_URL:=postgresql://root:password@localhost:5432/grpc?sslmode=disable

createdb:
	docker exec grpc-app-db createdb --username=root --owner=root grpc

dropdb:
	docker exec grpc-app-db dropdb grpc

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

# rollback last migration
migratedown1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down 1

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/devder/grpc-b/db/sqlc Store

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb dropdb migrateup migratedown sqlc server mock test migratedown1

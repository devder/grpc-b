createdb:
	docker exec grpc-app-db createdb --username=root --owner=root grpc

dropdb:
	docker exec grpc-app-db dropdb grpc

migrateup:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/grpc?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/grpc?sslmode=disable" -verbose down

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/devder/grpc-b/db/sqlc Store

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb dropdb migrateup migratedown sqlc server mock test

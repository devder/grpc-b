createdb:
	docker exec grpc-app-db createdb --username=root --owner=root grpc

dropdb:
	docker exec grpc-app-db dropdb grpc

migrateup:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/grpc?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5432/grpc?sslmode=disable" -verbose down

.PHONY: createdb dropdb migrateup migratedown

createdb:
	docker exec -it grpc-app-db createdb --username=root --owner=root grpc

dropdb:
	docker exec -it grpc-app-db dropdb grpc

.PHONY: createdb dropdb
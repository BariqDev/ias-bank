postgres: 
		docker run --name postgres-17-alpine -p 5433:5432  -e POSTGRES_USER=root -e  POSTGRES_PASSWORD=secret -d postgres:17.1-alpine3.20
createdb:
		docker exec -it postgres-17-alpine createdb --username=root --owner=root ias_bank
dropdb:
		docker exec -it postgres-17-alpine dropdb ias_bank
migrateup:
		migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/ias_bank?sslmode=disable" -verbose up
migrateup1:
		migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/ias_bank?sslmode=disable" -verbose up 1
		
migratedown:
		migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/ias_bank?sslmode=disable" -verbose down
migratedown1:
		migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/ias_bank?sslmode=disable" -verbose down 1
sqlc:
		sqlc generate
server:
		go run main.go

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock   server
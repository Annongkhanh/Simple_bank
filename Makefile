postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
simple_bank:
	docker run --name simplebank --network bank-network -e GIN_MODE=release -e  DB_SOURCE="postgresql://root:secret@postgres12:5432/simple_bank?sslmode=disable" -p 8080:8080 simplebank:latest
createdb:
	docker exec -it postgres12 createdb --user=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
pullimage:
	docker pull postgres:12-alpine
sqlc:
	sqlc generate
init:
	make postgres
	make createdb
	make migrateup
removeall:
	sudo docker rm -f $(sudo docker ps -aq)
test:
	go test -v -cover ./...
server:
	go run main.go 
mock:
	mockgen -destination db/mock/store.go -package mockdb  github.com/Annongkhanh/Go_example/db/sqlc Store
.PHONY: createdb dropdb postgres migrateup migratedown pull sqlc test server mock migrateup1 migratedown1

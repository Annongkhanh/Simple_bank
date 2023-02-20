postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --user=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
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
.PHONY: createdb dropdb postgres migrateup migratedown pull sqlc test server mock

package main

import (
	"database/sql"
	"log"

	"github.com/Annongkhanh/Go_example/api"
	db "github.com/Annongkhanh/Go_example/db/sqlc"
	"github.com/Annongkhanh/Go_example/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can not load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Can not connect to database: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Can not initialize server: ", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Can not start server: ", err)
	}

}

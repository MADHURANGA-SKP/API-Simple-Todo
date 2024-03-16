package main

import (
	"database/sql"
	"log"
	"simpletodo/api"
	db "simpletodo/db/sqlc"

	util "simpletodo/util"

	_ "github.com/lib/pq"
)

const (
	
)

func main (){
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	store := db.NewStore(conn)
	server, err:= api.NewServer(config, *store)
	if err != nil{
		log.Fatal("cannot create server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:",err)
	}
}
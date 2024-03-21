package main

import (
	"database/sql"
	"log"
	"net"
	"simpletodo/api"
	db "simpletodo/db/sqlc"
	"simpletodo/gapi"
	"simpletodo/pb"
	"time"

	util "simpletodo/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	
)

func main (){
	router := gin.Default()

    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://*", "https://*", "*", "https://testnet.bethelnet.io"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot start server", err)
	}

	store := db.NewStore(conn)
	runGrpcServer(config, *store)
	// runGinServer(config, *store)

}

func runGrpcServer(config util.Config, store db.Store){
	server, err := gapi.NewServer(config,store)
	if err != nil {
		log.Fatal("cannont create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpletodoServer(grpcServer,server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("Start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server")
	}
}

func runGinServer(config util.Config, store db.Store){
	server, err:= api.NewServer(config, store)
	if err != nil{
		log.Fatal("cannot create server", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:",err)
	}
}
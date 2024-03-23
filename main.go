package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"simpletodo/api"
	db "simpletodo/db/sqlc"
	_ "simpletodo/doc/statik"
	"simpletodo/gapi"
	"simpletodo/pb"
	util "simpletodo/util"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	 
	
	go runGatwayServer(config, *store)
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
		log.Fatal("cannot create listener",err)
	}

	log.Printf("Start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start grpc server",err)
	}
}

func runGatwayServer(config util.Config, store db.Store){
	server, err := gapi.NewServer(config,store)
	if err != nil {
		log.Fatal("cannont create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpletodoHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot registerhandler server:",err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil{
		log.Fatal("cannot create statik fs",err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/",swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:",err)
	}

	log.Printf("Start http gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start http gateway server:",err)
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
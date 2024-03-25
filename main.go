package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"simpletodo/api"
	db "simpletodo/db/sqlc"
	_ "simpletodo/doc/statik"
	"simpletodo/gapi"
	"simpletodo/pb"
	util "simpletodo/util"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
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
		log.Fatal().Msg("cannot connect to db:")
	}

	if config.Enviornment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}	

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}

	runDBMigeations(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	 
	
	go runGatwayServer(config, *store)
	runGrpcServer(config, *store)
	// runGinServer(config, *store)

}

func runDBMigeations(migrationURL string, dbSource string){
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil { 
		log.Fatal().Msg("cannot create ne migrate instance:")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange{
		log.Fatal().Msg("failed to run migrate up:")
	} 

	log.Info().Msg("db migrated succesfully")
}
 
func runGrpcServer(config util.Config, store db.Store){
	server, err := gapi.NewServer(config,store)
	if err != nil {
		log.Fatal().Msg("cannont create server:")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)

	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpletodoServer(grpcServer,server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	log.Info().Msgf("Start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start grpc server")
	}
}




func runGatwayServer(config util.Config, store db.Store){
	server, err := gapi.NewServer(config,store)
	if err != nil {
		log.Fatal().Msg("cannont create server:")
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
		log.Fatal().Msg("cannot registerhandler server:")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil{
		log.Fatal().Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/",swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener:")
	}

	log.Info().Msgf("Start http gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal().Msg("cannot start http gateway server:")
	}
}

func runGinServer(config util.Config, store db.Store){
	server, err:= api.NewServer(config, store)
	if err != nil{
		log.Fatal().Msg("cannot create server")
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server:")
	}
}
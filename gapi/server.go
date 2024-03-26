package gapi

import (
	"fmt"
	db "simpletodo/db/sqlc"
	"simpletodo/pb"
	"simpletodo/token"
	util "simpletodo/util"
	"simpletodo/worker"
)

//server serves gRPC requests
type Server struct {
	pb.UnimplementedSimpletodoServer
    config     util.Config
    store      db.Store
    tokenMaker token.Maker
	taskDistributor worker.TaskDistributor
}

//NewServer creates a gRPC server and setup routing
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
    tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
    if err != nil {
        return nil, fmt.Errorf("cannot create token maker: %w", err)
    }

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
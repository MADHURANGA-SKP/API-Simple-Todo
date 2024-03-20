package gapi

import (
	"fmt"
	db "simpletodo/db/sqlc"
	"simpletodo/pb"
	"simpletodo/token"
	util "simpletodo/util"
)

//server serves gRPC requests
type Server struct {
	pb.UnimplementedSimpletodoServer
    config     util.Config
    store      db.Store
    tokenMaker token.Maker
}

//NewServer creates a gRPC server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
    tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
    if err != nil {
        return nil, fmt.Errorf("cannot create token maker: %w", err)
    }

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
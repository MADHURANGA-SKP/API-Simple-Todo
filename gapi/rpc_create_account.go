package gapi

import (
	"context"
	db "simpletodo/db/sqlc"
	"simpletodo/pb"
	util "simpletodo/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateAccount(ctx context.Context,req *pb.CreateAccountRequest) (*pb.CreateAccountResult, error) {
	hashPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}
	
	arg := db.CreateAccountsParams{     
		FirstName: req.FirstName,
		LastName: req.LastName,
		UserName: req.UserName,
		Password: hashPassword,
	}

	account, err := server.store.CreateAccount(ctx,arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	
	rsp := &pb.CreateAccountResult{
		Account: convertAccount(account.Account),
	}

	return rsp, nil

}
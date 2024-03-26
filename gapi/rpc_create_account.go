package gapi

import (
	"context"
	db "simpletodo/db/sqlc"
	"simpletodo/pb"
	util "simpletodo/util"
	"simpletodo/val"
	"simpletodo/worker"
	"time"

	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateAccount(ctx context.Context,req *pb.CreateAccountRequest) (*pb.CreateAccountResult, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}
	
	arg := db.CreateAccountsParams{     
		FirstName: req.GetFirstName(),
		LastName: req.GetLastName(),
		Email: req.GetEmail(),
		UserName: req.GetUserName(),
		Password: hashPassword,
	}

	account, err := server.store.CreateAccount(ctx,arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}
	taskPayload := &worker.PayloadSendVerifyEmail{
		UserName: account.Account.UserName,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10*time.Second),
		asynq.Queue("critical"),
	}
	err =  server.taskDistributor.DistributeTaskSendVerifyEmail(ctx,taskPayload, opts...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to distribute task to send verifications")
	}


	rsp := &pb.CreateAccountResult{
		Account: convertAccount(account.Account),
	}

	return rsp, nil

}

func validateCreateUserRequest(req *pb.CreateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation){
	if err := val.ValidateFirstname(req.GetFirstName()); err != nil {
		violations = append(violations, fieldViolation("first_name", err))
  	}
	
	if err := val.ValidateLastname(req.GetLastName()); err != nil {
		violations = append(violations, fieldViolation("last_name", err))
  	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
   	}

	if err := val.ValidateUsername(req.GetUserName()); err != nil {
		 violations = append(violations, fieldViolation("user_name", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
   	}

	return violations
}
package gapi

import (
	"context"
	"database/sql"
	db "simpletodo/db/sqlc"
	"simpletodo/pb"
	util "simpletodo/util"
	"simpletodo/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateAccount(ctx context.Context,req *pb.UpdateAccountRequest) (*pb.UpdateAccountResult, error) {
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	
	arg := db.UpdateAccountsParams{     
		
		FirstName: sql.NullString{
			String: req.GetFirstName(),
			Valid: req.FirstName != nil,
		},
		LastName: sql.NullString{
			String: req.GetLastName(),
			Valid: req.LastName != nil,
		},
		UserName: req.GetUserName(),
		Email: sql.NullString{
			String:  req.GetEmail(),
			Valid: req.Email != nil,
		},

	}

	if req.Password != nil {
		hashPassword, err := util.HashPassword(req.GetPassword())
			if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		arg.Password = sql.NullString{
			String: hashPassword,
			Valid: true,
		}
	}

	account, err := server.store.UpdateAccount(ctx,arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	
	rsp := &pb.UpdateAccountResult{
		Account: convertAccount(account.Account),
	}

	return rsp, nil

}

func validateUpdateUserRequest(req *pb.UpdateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation){
	if req.FirstName != nil {
		if err := val.ValidateFirstname(req.GetFirstName()); err != nil {
			violations = append(violations, fieldViolation("first_name", err))
		  }
	}
	
	if req.LastName != nil {
		if err := val.ValidateLastname(req.GetLastName()); err != nil {
			violations = append(violations, fieldViolation("last_name", err))
		  }
	}

	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		   }
	}

	if err := val.ValidateUsername(req.GetUserName()); err != nil {
			violations = append(violations, fieldViolation("user_name", err))
	}
	

	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	return violations
}
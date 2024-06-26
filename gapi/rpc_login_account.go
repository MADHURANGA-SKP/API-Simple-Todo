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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginAccount(ctx context.Context, req *pb.LoginAccountRequest) (*pb.LoginAccountResult, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	
	account , err := server.store.GetAccount(ctx, db.GetAccountsParams{UserName: req.GetUserName()})
	if err != nil {
		if err == sql.ErrNoRows{
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to found user")
	}

	HashPassword, err := util.HashPassword(req.Password)
	
	err = util.CheckPassword(req.Password, HashPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		account.Account.UserName, 
		req.Password,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create access token")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		account.Account.UserName, 
		req.Password,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create refresh token")
	}

	mtd := server.extractMetadata(ctx)

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:     refreshPayload.ID,     
		AccountID: account.Account.ID,
		RefreshToken: refreshToken,
		UserAgent:    mtd.UserAgent,
		ClientIp:     mtd.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create session")

	}

	rsp := &pb.LoginAccountResult{
		Account: convertAccount(account.Account),
		SessionId: session.ID.String(),
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return rsp, nil
}

func validateLoginUserRequest(req *pb.LoginAccountRequest) (violations []*errdetails.BadRequest_FieldViolation){
	if err := val.ValidateUsername(req.GetUserName()); err != nil {
		 violations = append(violations, fieldViolation("user_name", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
   	}

	return violations
}

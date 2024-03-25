package gapi

import (
	db "simpletodo/db/sqlc"
	"simpletodo/pb"
)

func convertAccount(account db.Account) *pb.Account{
	return &pb.Account{
		FirstName: account.FirstName,
		LastName: account.LastName,
		Email: account.Email,
		UserName: account.UserName,
		Password: account.Password,
	}
}
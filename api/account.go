package api

import (
	"database/sql"
	"errors"
	"time"

	"net/http"

	db "simpletodo/db/sqlc"
	"simpletodo/token"

	util "simpletodo/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// createAccountRequest contains the input parameters for create an account in account table
type createAccountRequest struct{
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	UserName  string `json:"user_name" binding:"required,alphanum,min=8"`
}

type AccountResult struct{
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	UserName  string `json:"user_name" binding:"required,alphanum,min=8"`
	
}


func newAccountResult(user db.Account) AccountResult {
	return AccountResult{
			FirstName: user.FirstName,
			LastName: user.LastName,
			UserName: user.UserName,
	}
	
}

func (server *Server) CreateAccount(ctx *gin.Context){
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// // Check if authorization payload exists in the context
    // authPayload, exists := ctx.Get(authorizationPayloadKey)
    // if !exists {
    //     err := errors.New("authorization payload is missing")
    //     ctx.JSON(http.StatusUnauthorized, errorResponse(err))
    //     return
    // }

    // // Assert the authorization payload to the correct type
    // payload, ok := authPayload.(*token.Payload)
    // if !ok {
    //     err := errors.New("authorization payload is invalid")
    //     ctx.JSON(http.StatusUnauthorized, errorResponse(err))
    //     return
    // }


	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
    ctx.JSON(http.StatusInternalServerError, errorResponse(err))
    return
	}
	
	arg := db.CreateAccountsParams{     
		FirstName: req.FirstName,
		LastName: req.LastName,
		UserName: req.UserName,
		Password: hashPassword,
	}

	user, err := server.store.CreateAccount(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	
	rsp := newAccountResult(user.Account)
	ctx.JSON(http.StatusOK, rsp)
}

// GetAccountsRequest contains the input parameters and process to request for data from account table
type GetAccountsRequest struct {
	UserName  string `json:"user_name"`
}

func (server *Server) GetAccount(ctx *gin.Context){
	var req GetAccountsRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetAccountsParams{UserName: req.UserName}

	account, err := server.store.GetAccount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Account.UserName != authPayload.Username{
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
	}


	ctx.JSON(http.StatusOK, account)
}

// type listAccountRequest struct {
// 	PageID   int32 `form:"page_id" binding:"required,min=1"`
// 	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
// }

// func (server *Server) listAccounts(ctx *gin.Context) {
// 	var req listAccountRequest
// 	if err := ctx.ShouldBindQuery(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	// authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
// 	arg := db.ListAccountParams{
// 		// UserName: authPayload.Username,
// 		Limit:  req.PageSize,
// 		Offset: (req.PageID - 1) * req.PageSize,
// 	}

// 	accounts, err := server.store.ListAccount(ctx, arg)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, accounts)
// }

type loginAccountRequest struct {
    UserName string `json:"user_name" binding:"required,alphanum"`
    Password  string `json:"password" validate:"required min=6"`
}

type loginAccountResult struct {
	SessionID	uuid.UUID `json:"session_id"`
    AccessToken string       `json:"access_token"`
	AccesssTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshToken string `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
	Account AccountResult `json:"account"`
}




func (server *Server) LoginAccount (ctx *gin.Context){
	var req loginAccountRequest
	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account , err := server.store.GetAccount(ctx, db.GetAccountsParams{UserName: req.UserName})
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	HashPassword, err := util.HashPassword(req.Password)
	
	err = util.CheckPassword(req.Password, HashPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		account.Account.UserName, 
		req.Password,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		account.Account.UserName, 
		req.Password,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:     refreshPayload.ID,     
		AccountID: account.Account.ID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginAccountResult{
		SessionID: session.ID,
		AccessToken: accessToken,
		AccesssTokenExpiresAt: accessPayload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		Account: newAccountResult(account.Account),
	}

	ctx.JSON(http.StatusOK, rsp)
} 
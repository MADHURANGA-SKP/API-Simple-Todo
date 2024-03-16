package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	db "simpletodo/db/sqlc"
	"simpletodo/token"

	"github.com/gin-gonic/gin"
)

//createTodoRequest contains the input parameters for create an Todo
type createTodoRequest struct{
	Username string `json:"username"`
	Title    string      `json:"title"`
	Time     string      `json:"time"`
	Date     string      `json:"date"`
	Complete string `json:"complete"`
}

func (server *Server) CreateTodo(ctx *gin.Context){
	var req createTodoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload, ok := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if !ok {
        err := errors.New("authorization payload is invalid")
        ctx.JSON(http.StatusUnauthorized, errorResponse(err))
        return
    }

    if req.Username != authPayload.Username {
        err := errors.New("account doesn't belong to the authenticated user")
        ctx.JSON(http.StatusUnauthorized, errorResponse(err))
        return
    }

    _, valid := server.validTodo(ctx, req.Username)
    if !valid {
        return
    }
	arg := db.CreateTodosParams{
		Title: req.Title,
		Time: req.Time,
		Date: req.Date,
		Complete: req.Complete,
	}

	Todo, err := server.store.CreateTodo(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, Todo)
}

// GettodosRequest contains the input parameters and process to request for data from account table
type GetTodosRequest struct {
	AccountID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetTodo(ctx *gin.Context){
	var req GetTodosRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.GetTodoParams{AccountID: req.AccountID}

	account, err := server.store.GetTodo(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}


//UpdateTodosRequest contains the input parameters for update an todo
type UpdateTodosRequest struct {
	Title    string      `json:"title"`
    Time     string      `json:"time"`
    Date     string      `json:"date"`
    Complete string `json:"complete"`

}

func (server *Server) UpdateTodo(ctx *gin.Context){
	var req UpdateTodosRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateTodosParams{}

	Todo, err := server.store.UpdateTodo(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, Todo)
}

// DeleteTodoRequest contains the input parameters for deleting an todo
type DeleteTodoRequest struct {
    ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) DeleteTodo(ctx *gin.Context) {
    var req DeleteTodoRequest
    if err := ctx.ShouldBindUri(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    err := server.store.DeleteTodo(ctx, req.ID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Owner deleted successfully"})
}


//ListTodoRequest contains the input parameters for deleting an todo
type listTodoRequest struct {
    PageID   int32 `form:"page_id" binding:"required,min=1"`
    PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) ListTodo(ctx *gin.Context) {
    var req listTodoRequest
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    arg := db.ListTodoParams{
        Limit:  req.PageSize,
        Offset: (req.PageID - 1) * req.PageSize,
    }

    owners, err := server.store.ListTodo(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, owners)
}
	

func (server *Server) validTodo(ctx *gin.Context,  username string) (db.Account, bool) {
    arg := db.GetAccountsParams{UserName: username}  
	account, err := server.store.GetAccount(ctx, arg)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return db.Account{}, false
        }

        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return db.Account{}, false
    }

    if account.Account.UserName != username {
        err := fmt.Errorf("account username mismatch: %s vs ", username)
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return db.Account{}, false
    }

    return account.Account, true
}


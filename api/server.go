package api

import (
	"fmt"
	db "simpletodo/db/sqlc"

	"simpletodo/token"

	util "simpletodo/util"

	"github.com/gin-gonic/gin"
)

//server serves hhtp requests
type Server struct {
    config     util.Config
    store      db.Store
    tokenMaker token.Maker
    router     *gin.Engine
}

//NewServer creates a http server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
    tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
    if err != nil {
        return nil, fmt.Errorf("cannot create token maker: %w", err)
    }

    server := &Server{
        config:     config,
        store:      store,
        tokenMaker: tokenMaker,
    }
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter(){
	router := gin.Default()

	router.POST("/account/create", server.CreateAccount)
	router.POST("/account/login", server.LoginAccount)
	router.GET("/account/get/:id", server.GetAccount)
	// router.PUT("/account/update/:id", server.UpdateAccount)
	// router.GET("/account/list", server.listAccounts)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))


	authRoutes.POST("/todo/create", server.CreateTodo)
	authRoutes.PUT("/todo/update/:id", server.UpdateTodo)
	authRoutes.DELETE("/todo/delete/:id", server.DeleteTodo)
	authRoutes.GET("/todo/get/:id", server.GetTodo)
	server.router = router
}

//start runs the http server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H{
	return gin.H{"error": err.Error()}
}
package api

import (
	"fmt"

	db "github.com/devder/grpc-b/db/sqlc"
	"github.com/devder/grpc-b/token"
	"github.com/devder/grpc-b/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

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

	// register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", util.ValidCurrency)
	}

	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	err := router.SetTrustedProxies(nil)

	if err != nil {
		fmt.Printf("Failed to set trusted proxies: %v", err)
	}

	// users endpoints
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/refresh", server.renewAccessToken)

	authRoute := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// accounts endpoints
	authRoute.POST("/accounts", server.createAccount)
	authRoute.GET("/accounts/:id", server.getAccount)
	authRoute.GET("/accounts", server.getAccounts)

	// transfers endpoints
	authRoute.POST("/transfers", server.createTransfer)
	server.router = router
}

// Run the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

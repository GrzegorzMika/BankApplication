package api

import (
	"fmt"
	"log"

	"BankApplication/internal/db"
	"BankApplication/internal/token"
	"BankApplication/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for banking application.
type Server struct {
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		router:     gin.Default(),
		config:     config,
	}

	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := val.RegisterValidation("currency", ValidCurrency)
		if err != nil {
			log.Println("Failed to register validator: ", err)
		}
	}
	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	s.router.POST("/users", s.createUser)
	s.router.POST("/users/login", s.loginUser)
	s.router.POST("/token/renew_access", s.refreshAccessToken)

	authRoutes := s.router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.GET("/accounts", s.listAccount)

	authRoutes.POST("/transfers", s.createTransfer)
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

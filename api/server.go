package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/valverdethiago/trading-api/db/sqlc"
)

// Server serves HTTP requests for the stock trading REST API
type Server struct {
	queries db.Querier
	router  *gin.Engine
}

// NewServer queries a new HTTP Server for the REST API
func NewServer(queries db.Querier) *Server {
	server := &Server{
		queries: queries,
		router:  gin.Default(),
	}
	accountController := NewAccountController(server.queries)
	accountController.setupRoutes(server.router)
	addressController := NewAddressController(server.queries)
	addressController.setupRoutes(server.router)
	tradeController := NewTradeController(server.queries)
	tradeController.setupRoutes(server.router)
	return server
}

// Start runs the HTTP Server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

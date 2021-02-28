package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/valverdethiago/trading-api/db/sqlc"
)

const accountsPath = "/accounts"

type createAccountRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required"`
	Address  *address `json:"address" binding:"omitempty"`
}

type address struct {
	Name    string `json:"name" binding:"required"`
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state" binding:"required"`
	Zipcode string `json:"zipcode" binding:"required"`
}

// AccountController controller for accounts object
type AccountController struct {
	queries *db.Queries
}

// NewAccountController builds a new instance of account controller
func NewAccountController(queries *db.Queries) *AccountController {
	return &AccountController{
		queries: queries,
	}

}

func (controller *AccountController) setupRoutes(router *gin.Engine) {
	router.POST(accountsPath, controller.createAccount)
	router.GET(accountsPath, controller.listAccounts)
}

func (controller *AccountController) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	message := fmt.Sprintf("Received name %s", req.Username)
	log.Printf(message)
	ctx.JSON(http.StatusAccepted, gin.H{"message": message})
	return
}

func (controller *AccountController) listAccounts(ctx *gin.Context) {
	accounts, err := controller.queries.ListAccounts(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if accounts == nil || len(accounts) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "No accounts found"})
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

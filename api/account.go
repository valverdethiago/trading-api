package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/valverdethiago/trading-api/db/sqlc"
	"github.com/valverdethiago/trading-api/service"
)

const (
	accountsPath     = "/accounts"
	accountsPathByID = "/accounts/:id"
)

type accountIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

// CreateAccountRequest json request to create account
type CreateAccountRequest struct {
	Username string          `json:"username" binding:"required"`
	Email    string          `json:"email" binding:"required"`
	Address  *AddressRequest `json:"address" binding:"omitempty"`
}

// AddressRequest json request to create address
type AddressRequest struct {
	Name    string `json:"name" binding:"required"`
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state" binding:"required"`
	Zipcode string `json:"zipcode" binding:"required"`
}

// AccountController controller for accounts object
type AccountController struct {
	service *service.AccountService
}

// NewAccountController builds a new instance of account controller
func NewAccountController(queries db.Querier) *AccountController {
	return &AccountController{
		service: service.NewAccountService(queries),
	}

}

func (controller *AccountController) setupRoutes(router *gin.Engine) {
	router.POST(accountsPath, controller.createAccount)
	router.GET(accountsPath, controller.listAccounts)
	router.GET(accountsPathByID, controller.findAccountByID)
}

func (controller *AccountController) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var account = db.Account{
		Username: req.Username,
		Email:    req.Email,
	}
	var address *db.Address
	if req.Address != nil {
		address = &db.Address{
			Name:    req.Address.Name,
			Street:  req.Address.Street,
			City:    req.Address.City,
			State:   db.State(req.Address.State),
			Zipcode: req.Address.Zipcode,
		}
	}
	dbAccount, _, err := controller.service.CreateAccount(account, address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dbAccount)
	return
}

func (controller *AccountController) listAccounts(ctx *gin.Context) {
	accounts, err := controller.service.ListAccounts()
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if accounts == nil || len(accounts) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "No accounts found"})
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

type getAccountRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (controller *AccountController) findAccountByID(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	uuid, err := parseUUID(req.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := controller.service.GetAccountByID(uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "No account found for this id"})
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	ctx.JSON(http.StatusOK, account)

}

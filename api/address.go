package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/valverdethiago/trading-api/db/sqlc"
	"github.com/valverdethiago/trading-api/service"
)

const (
	addressPath = "/accounts/:id/address"
)

type persistAddressRequest struct {
	Name    string `json:"name" binding:"required"`
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state" binding:"required"`
	Zipcode string `json:"zipcode" binding:"required"`
}

// AddressController controller for address object
type AddressController struct {
	service *service.AddressService
}

// NewAddressController builds a new instance of account controller
func NewAddressController(queries *db.Queries) *AddressController {
	accountService := service.NewAccountService(queries)
	return &AddressController{
		service: service.NewAddressService(
			queries, accountService,
		),
	}

}

func (controller *AddressController) setupRoutes(router *gin.Engine) {
	router.POST(addressPath, controller.updateAddressForAccount)
	router.GET(addressPath, controller.getAddressByAccountID)
	router.PUT(addressPath, controller.createAddressForAccount)
}

func (controller *AddressController) createAddressForAccount(ctx *gin.Context) {
	idReq, err := getAccountIDRequest(ctx)
	if err != nil {
		return
	}
	req, err := getAddressRequest(ctx)
	if err != nil {
		return
	}
	var address = db.Address{
		Name:    req.Name,
		Street:  req.Street,
		City:    req.City,
		State:   db.State(req.State),
		Zipcode: req.Zipcode,
	}
	dbAddress, err := controller.service.CreateAddressForAccount(idReq.ID, address)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dbAddress)
	return
}

func (controller *AddressController) updateAddressForAccount(ctx *gin.Context) {
	idReq, err := getAccountIDRequest(ctx)
	if err != nil {
		return
	}
	req, err := getAddressRequest(ctx)
	if err != nil {
		return
	}
	var address = db.Address{
		Name:    req.Name,
		Street:  req.Street,
		City:    req.City,
		State:   db.State(req.State),
		Zipcode: req.Zipcode,
	}
	dbAddress, err := controller.service.UpdateAddressForAccount(idReq.ID, address)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dbAddress)
	return
}

func (controller *AddressController) getAddressByAccountID(ctx *gin.Context) {
	idReq, err := getAccountIDRequest(ctx)
	if err != nil {
		return
	}
	dbAddress, err := controller.service.GetAddressByAccountID(idReq.ID)
	if err != nil && err == sql.ErrNoRows {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "The account doesn't have an address yet."})
		return
	}
	ctx.JSON(http.StatusOK, dbAddress)
	return
}

func getAccountIDRequest(ctx *gin.Context) (accountIDRequest, error) {
	var idReq accountIDRequest
	var err error
	if err = ctx.ShouldBindUri(&idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	return idReq, err
}

func getAddressRequest(ctx *gin.Context) (addressRequest, error) {
	var req addressRequest
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	return req, err
}

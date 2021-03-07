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

// AddressController controller for address object
type AddressController struct {
	service *service.AddressService
}

// NewAddressController builds a new instance of account controller
func NewAddressController(queries db.Querier) *AddressController {
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
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	uuid, err := parseUUID(idReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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
	dbAddress, err := controller.service.CreateAddressForAccount(uuid, address)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusConflict, errorResponse(err))
		}
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
	uuid, err := parseUUID(idReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
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
	dbAddress, err := controller.service.UpdateAddressForAccount(uuid, address)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	ctx.JSON(http.StatusCreated, dbAddress)
	return
}

func (controller *AddressController) getAddressByAccountID(ctx *gin.Context) {
	idReq, err := getAccountIDRequest(ctx)
	if err != nil {
		return
	}
	uuid, err := parseUUID(idReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	dbAddress, err := controller.service.GetAddressByAccountID(uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
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

func getAddressRequest(ctx *gin.Context) (AddressRequest, error) {
	var req AddressRequest
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	return req, err
}

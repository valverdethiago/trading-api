package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/valverdethiago/trading-api/db/sqlc"
	"github.com/valverdethiago/trading-api/service"
)

const (
	tradesPath     = "/accounts/:id/trades"
	tradesPathByID = "/accounts/:id/trades/:tradeID"
)

type tradeRequest struct {
	Symbol   string       `json:"symbol" binding:"required"`
	Quantity int64        `json:"quantity" binding:"required,min=1"`
	Side     db.TradeSide `json:"side" binding:"required"`
	Price    float64      `json:"price" binding:"required,gt=0"`
}

type tradeIDRequest struct {
	ID string `uri:"tradeID" binding:"required"`
}

// TradeController controller for trades object
type TradeController struct {
	service *service.TradeService
}

// NewTradeController builds a new intance of trade controller
func NewTradeController(queries db.Querier) *TradeController {
	accountService := service.NewAccountService(queries)
	return &TradeController{
		service: service.NewTradeService(queries, accountService),
	}
}

func (controller *TradeController) setupRoutes(router *gin.Engine) {
	router.POST(tradesPath, controller.createTrade)
	router.GET(tradesPath, controller.listTradesByAccount)
	router.GET(tradesPathByID, controller.getTradeByIDAndAccountID)
	router.DELETE(tradesPathByID, controller.cancelTradeByIDAndAccountID)
}

func (controller *TradeController) createTrade(ctx *gin.Context) {
	idReq, err := getAccountIDRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	req, err := getTradeRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	accountUUID, err := parseUUID(idReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	trade := db.Trade{
		Symbol:   req.Symbol,
		Side:     req.Side,
		Price:    req.Price,
		Quantity: req.Quantity,
	}
	dbTrade, err := controller.service.CreateTrade(trade, accountUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	ctx.JSON(http.StatusCreated, dbTrade)
}

func (controller *TradeController) listTradesByAccount(ctx *gin.Context) {
	idReq, err := getAccountIDRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accountUUID, err := parseUUID(idReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	dbTrades, err := controller.service.ListTradesByAccount(accountUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		return
	}
	ctx.JSON(http.StatusOK, dbTrades)
}

func (controller *TradeController) getTradeByIDAndAccountID(ctx *gin.Context) {
	accountIDReq, err := getAccountIDRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	tradeIDReq, err := getTradeIDRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	accountUUID, err := parseUUID(accountIDReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	tradeUUID, err := parseUUID(tradeIDReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	dbTrade, err := controller.service.FindByIDAndAccountID(tradeUUID, accountUUID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dbTrade)
}

func (controller *TradeController) cancelTradeByIDAndAccountID(ctx *gin.Context) {
	accountIDReq, err := getAccountIDRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	tradeIDReq, err := getTradeIDRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accountUUID, err := parseUUID(accountIDReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	tradeUUID, err := parseUUID(tradeIDReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	dbTrade, err := controller.service.CancelTradeByIDAndAccountID(tradeUUID, accountUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		} else {
			ctx.JSON(http.StatusConflict, errorResponse(err))
		}
		return
	}
	ctx.JSON(http.StatusAccepted, dbTrade)
}

func getTradeRequest(ctx *gin.Context) (tradeRequest, error) {
	var req tradeRequest
	var err error
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	return req, err
}

func getTradeIDRequest(ctx *gin.Context) (tradeIDRequest, error) {
	var req tradeIDRequest
	var err error
	if err = ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	return req, err
}

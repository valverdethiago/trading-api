package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/valverdethiago/trading-api/db/sqlc"
	"github.com/valverdethiago/trading-api/db/store"
	"github.com/valverdethiago/trading-api/service"
)

const (
	tradesPath     = "/accounts/:id/trades"
	tradesPathByID = "/accounts/:id/trades/:tradeID"
)

type tradeRequest struct {
	Symbol   string         `json:"symbol"`
	Quantity int64          `json:"quantity"`
	Side     db.TradeSide   `json:"side"`
	Price    float64        `json:"price"`
	Status   db.TradeStatus `json:"status"`
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
	accountStore := store.NewDbAccountStore(queries)
	addressStore := store.NewDbAddressStore(queries)
	tradeStore := store.NewDbTradeStore(queries)
	accountService := service.NewAccountService(accountStore, addressStore)
	return &TradeController{
		service: service.NewTradeService(tradeStore, accountService),
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
	trade := db.Trade{
		Symbol:   req.Symbol,
		Side:     req.Side,
		Price:    req.Price,
		Quantity: req.Quantity,
	}
	dbTrade, err := controller.service.CreateTrade(trade, idReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dbTrade)
}

func (controller *TradeController) listTradesByAccount(ctx *gin.Context) {
	idReq, err := getAccountIDRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	dbTrades, err := controller.service.ListTradesByAccount(idReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dbTrades)
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
	dbTrade, err := controller.service.FindByIDAndAccountID(tradeIDReq.ID, accountIDReq.ID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dbTrade)
}

func (controller *TradeController) cancelTradeByIDAndAccountID(ctx *gin.Context) {
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
	dbTrade, err := controller.service.CancelTradeByIDAndAccountID(tradeIDReq.ID, accountIDReq.ID)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, dbTrade)
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
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}
	return req, err
}

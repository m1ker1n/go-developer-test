package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/m1ker1n/go-developer-test/internal/models"
	"github.com/m1ker1n/go-developer-test/internal/services"
	"github.com/m1ker1n/go-developer-test/internal/storage"
	"github.com/shopspring/decimal"
	"net/http"
)

type walletService interface {
	CreateWallet(ctx context.Context) (models.Wallet, error)
	GetWallet(ctx context.Context, walletId uuid.UUID) (models.Wallet, error)
}

type transactionService interface {
	CreateTransaction(ctx context.Context, from, to uuid.UUID, amount decimal.Decimal) (models.Transaction, error)
	GetTransactions(ctx context.Context, walletId uuid.UUID) ([]models.Transaction, error)
}

type WalletHandler struct {
	walletService      walletService
	transactionService transactionService
}

func (w *WalletHandler) CreateWallet(ctx *gin.Context) {
	wallet, err := w.walletService.CreateWallet(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusCreated, wallet)
}

type SendMoneyArguments struct {
	To     string  `json:"to" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

func (w *WalletHandler) SendMoney(ctx *gin.Context) {
	from := ctx.Param("walletId")
	if from == "" {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "source wallet is not provided"})
		return
	}
	fromUuid, err := uuid.Parse(from)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "couldn't parse source wallet id"})
		return
	}

	var bodyArgs SendMoneyArguments
	if err := ctx.ShouldBindJSON(&bodyArgs); err != nil {
		//TODO: something went wrong with binding maybe must be other code
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	toUuid, err := uuid.Parse(bodyArgs.To)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "couldn't parse destination wallet id"})
		return
	}
	amount := decimal.NewFromFloat(bodyArgs.Amount)

	_, err = w.transactionService.CreateTransaction(ctx, fromUuid, toUuid, amount)
	if err != nil {
		if errors.Is(err, services.ErrTransactionWalletFromNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err})
			return
		}
		if errors.Is(err, services.ErrNotEnoughMoney) || errors.Is(err, services.ErrTransactionWalletToNotFound) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.Status(http.StatusOK)
}

func (w *WalletHandler) GetTransactionHistory(ctx *gin.Context) {
	from := ctx.Param("walletId")
	if from == "" {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "source wallet is not provided"})
		return
	}
	fromUuid, err := uuid.Parse(from)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "couldn't parse source wallet id"})
		return
	}

	transactions, err := w.transactionService.GetTransactions(ctx, fromUuid)
	if err != nil {
		if errors.Is(err, storage.ErrWalletNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}

func (w *WalletHandler) GetWalletInfo(ctx *gin.Context) {
	walletId := ctx.Param("walletId")
	if walletId == "" {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "source wallet is not provided"})
		return
	}
	walletIdUuid, err := uuid.Parse(walletId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "couldn't parse source wallet id"})
		return
	}

	wallet, err := w.walletService.GetWallet(ctx, walletIdUuid)
	if err != nil {
		if errors.Is(err, services.ErrWalletNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, wallet)
}

func NewWalletHandler(walletService walletService, transactionService transactionService) *WalletHandler {
	return &WalletHandler{walletService: walletService, transactionService: transactionService}
}

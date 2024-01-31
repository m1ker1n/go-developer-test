package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type walletHandler interface {
	CreateWallet(ctx *gin.Context)
	SendMoney(ctx *gin.Context)
	GetTransactionHistory(ctx *gin.Context)
	GetWalletInfo(ctx *gin.Context)
}

type Server struct {
	server http.Server
}

func New(addr string, walletHandler walletHandler) *Server {
	router := registerHandlers(walletHandler)
	return &Server{
		server: http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func registerHandlers(walletHandler walletHandler) *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	apiV1 := api.Group("/v1")
	apiV1Wallet := apiV1.Group("/wallet")
	{
		apiV1Wallet.POST("", walletHandler.CreateWallet)
		apiV1Wallet.POST("/:walletId/send", walletHandler.SendMoney)
		apiV1Wallet.GET("/:walletId/history", walletHandler.GetTransactionHistory)
		apiV1Wallet.GET("/:walletId", walletHandler.GetWalletInfo)
	}
	return router
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

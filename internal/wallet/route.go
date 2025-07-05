package wallet

import (
	"wallet-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, handler *WalletHandler) {
	walletGroup := r.Group("/api/v1/wallet") 
	{
		walletGroup.POST(":user_id", handler.CreateWallet)
		// walletGroup.GET("", handler.GetAllWallet)
		walletGroup.GET("/:user_id", handler.GetWalletByUserID)
		walletGroup.POST("/add_balance", middleware.Secured() , handler.AddBalance)
		walletGroup.POST("/deduct_balance", middleware.Secured() , handler.DeductBalance)
		// walletGroup.DELETE("/:id", handler.DeleteWallet)
	}
}
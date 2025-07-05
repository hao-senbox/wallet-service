package exchange

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *ExchangeHandler) {
    exchangeGroup := r.Group("/api/v1/exchange") 
    {
        exchangeGroup.POST("", handler.CreateExchangeRate)
        exchangeGroup.GET("", handler.GetAllExchangeRate)
        exchangeGroup.GET("/:id", handler.GetExchangeRate)
        exchangeGroup.PUT("/:id", handler.UpdateExchangeRate)
        exchangeGroup.DELETE("/:id", handler.DeleteExchangeRate)
    }
}
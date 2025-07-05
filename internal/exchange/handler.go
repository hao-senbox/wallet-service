package exchange

import (
	"net/http"
	"wallet-service/helper"

	"github.com/gin-gonic/gin"
)

type ExchangeHandler struct {
	exchangeService ExchangeService
}

func NewExchangeHandler(exchangeService ExchangeService) *ExchangeHandler {
	return &ExchangeHandler{
		exchangeService: exchangeService,
	}
}

func (h *ExchangeHandler) CreateExchangeRate(c *gin.Context) {
	
	var req CreateExchangeRateRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.exchangeService.CreateExchangeRate(c, &req)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}

func (h *ExchangeHandler) GetAllExchangeRate(c *gin.Context) {
	
	exchange, err := h.exchangeService.GetAllExchangeRate(c)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", exchange)

}

func (h *ExchangeHandler) GetExchangeRate(c *gin.Context) {
	
	id := c.Param("id")

	exchange, err := h.exchangeService.GetExchangeRate(c, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", exchange)

}

func (h *ExchangeHandler) UpdateExchangeRate(c *gin.Context) {
	
	id := c.Param("id")

	var req UpdateExchangeRateRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err := h.exchangeService.UpdateExchangeRate(c, id, &req)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}

func (h *ExchangeHandler) DeleteExchangeRate(c *gin.Context) {
	
	id := c.Param("id")

	err := h.exchangeService.DeleteExchangeRate(c, id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}
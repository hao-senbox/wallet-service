package wallet

import (
	"context"
	"fmt"
	"net/http"
	"wallet-service/helper"
	"wallet-service/pkg/constants"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	service WalletService
}

func NewWalletHandler(service WalletService) *WalletHandler {
	return &WalletHandler{
		service: service,
	}
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	
	user_id := c.Param("user_id")

	err := h.service.CreateWallet(c, user_id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}

func (h *WalletHandler) GetWalletByUserID(c *gin.Context) {
	
	user_id := c.Param("user_id")

	wallet, err := h.service.GetWalletByUserID(c, user_id)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", wallet)
	
}

func (h *WalletHandler) AddBalance(c *gin.Context) {
	
	var req AddBalanceRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	err := h.service.AddBalance(c, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}

func (h *WalletHandler) DeductBalance(c *gin.Context) {
	
	var req DeductBalanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	userID, exists := c.Get(constants.UserID)

	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	err := h.service.DeductBalance(ctx, &req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "success", nil)

}
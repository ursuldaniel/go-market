package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ursuldaniel/go-market/internal/domain/models"
)

func (s *Server) handleMakePurchase(c *gin.Context) {
	userId := c.MustGet("id").(int)

	productId, err := ParseId(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	quantity_ := c.Query("quantity")
	quantity, err := strconv.Atoi(quantity_)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.store.MakePurchase(userId, productId, quantity); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: "purchase successfully made"})
}

func (s *Server) handleGetUserPurchases(c *gin.Context) {
	userId := c.MustGet("id").(int)

	purchases, err := s.store.GetUserPurchases(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, purchases)
}

func (s *Server) handleGetProductPurchases(c *gin.Context) {
	productId, err := ParseId(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	purchases, err := s.store.GetProductPurchases(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, purchases)
}

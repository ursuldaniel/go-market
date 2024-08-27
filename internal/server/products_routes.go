package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ursuldaniel/go-market/internal/domain/models"
)

func (s *Server) handleAddProduct(c *gin.Context) {
	product := models.Product{}
	if err := c.ShouldBindBodyWithJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.store.AddProduct(product.Name, product.Description, product.Price, product.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: "product successfully added"})
}

func (s *Server) handleGetAllProducts(c *gin.Context) {
	products, err := s.store.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (s *Server) handleGetProductById(c *gin.Context) {
	id, err := ParseId(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	product, err := s.store.GetProductById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (s *Server) handleUpdateProduct(c *gin.Context) {
	id, err := ParseId(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	product := models.Product{}
	if err := c.ShouldBindBodyWithJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err = s.store.UpdateProduct(id, product.Name, product.Description, product.Price, product.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: "product successfully updated"})
}

func (s *Server) handleDeleteProduct(c *gin.Context) {
	id, err := ParseId(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.store.DeleteProduct(id); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: "product successfully deleted"})
}

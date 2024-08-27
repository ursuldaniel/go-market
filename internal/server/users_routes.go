package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ursuldaniel/go-market/internal/domain/models"
)

func (s *Server) handleRegisterUser(c *gin.Context) {
	user := models.User{}
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.validate.Struct(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.store.RegisterUser(user.Username, user.Password, user.Email); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: "user successfully created"})
}

func (s *Server) handleLoginUser(c *gin.Context) {
	loginUser := models.User{}
	if err := c.ShouldBindBodyWithJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.validate.Struct(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	id, err := s.store.LoginUser(loginUser.Username, loginUser.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	var token string

	if loginUser.Username == "admin" && loginUser.Password == "admin" {
		token, err = CreateAdminToken(id)
	} else {
		token, err = CreateUserToken(id)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.Response{Message: token})
}

func (s *Server) handleGetUserProfile(c *gin.Context) {
	id, err := ParseId(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	user, err := s.store.GetUserProfile(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) handleProfile(c *gin.Context) {
	id := c.MustGet("id").(int)

	user, err := s.store.GetUserProfile(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

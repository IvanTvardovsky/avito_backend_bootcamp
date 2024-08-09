package http

import (
	"avito_bootcamp/internal/entity"
	"avito_bootcamp/internal/usecase"
	"avito_bootcamp/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authRoutes struct {
	a usecase.Authorization
	l logger.Interface
}

func newAuthRoutes(handler *gin.Engine, l logger.Interface, a usecase.Authorization) {
	r := &authRoutes{
		a: a,
		l: l,
	}

	handler.GET("/dummyLogin", r.dummyLogin)
	handler.POST("/login", r.login)
	handler.POST("/register", r.register)
}

func (r *authRoutes) dummyLogin(c *gin.Context) {
	userType := c.Query("user_type")
	if userType == "" {
		errorResponse(c, r.l, http.StatusBadRequest, "http - dummyLogin: user_type is required", nil)
		return
	}

	token, err := r.a.DummyLogin(c.Request.Context(), userType)
	if err != nil {
		errorResponse(c, r.l, http.StatusBadRequest, "http - dummyLogin: failed to create dummy login token", err)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (r *authRoutes) login(c *gin.Context) {
	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, r.l, http.StatusBadRequest, "http - login: bad request body", err)
		return
	}

	token, err := r.a.Login(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, r.l, http.StatusUnauthorized, "http - login: invalid credentials", err)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (r *authRoutes) register(c *gin.Context) {
	var req entity.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, r.l, http.StatusBadRequest, "http - register: bad request body", err)
		return
	}

	user, err := r.a.Register(c.Request.Context(), req)
	if err != nil {
		errorResponse(c, r.l, http.StatusInternalServerError, "http - register: failed to register user", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": user.UserID})
}

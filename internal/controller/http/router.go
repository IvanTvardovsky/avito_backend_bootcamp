package http

import (
	"avito_bootcamp/internal/usecase"
	"avito_bootcamp/pkg/logger"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, l logger.Interface, JWTToken string, flatUC *usecase.FlatUseCase,
	houseUC *usecase.HouseUseCase, authUC *usecase.AuthUseCase) {

	handler.Use(gin.Logger())
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	newAuthRoutes(handler, l, authUC)

	handler.Use(AuthMiddleware(JWTToken))

	newFlatRoutes(handler, l, flatUC)
	newHouseRoutes(handler, l, houseUC)
}

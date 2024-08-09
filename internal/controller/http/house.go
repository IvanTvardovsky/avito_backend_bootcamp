package http

import (
	"avito_bootcamp/internal/entity"
	"avito_bootcamp/internal/usecase"
	"avito_bootcamp/pkg/logger"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type houseRoutes struct {
	h usecase.House
	l logger.Interface
}

func newHouseRoutes(handler *gin.Engine, l logger.Interface, h usecase.House) {
	r := &houseRoutes{
		h: h,
		l: l,
	}

	hand := handler.Group("/house/")
	{
		hand.POST("/create", RoleMiddleware([]string{"moderator"}), r.create)
		hand.GET("/:id", RoleMiddleware([]string{"moderator", "client"}), r.getFlats)
	}
}

type createHouseRequest struct {
	ID        int    `json:"id"`
	Address   string `json:"address"`
	Year      int    `json:"year"`
	Developer string `json:"developer"`
}

func (h *createHouseRequest) satisfiesTheConditions() bool {
	if len(h.Address) <= 0 || h.Year < 0 || h.ID < 1 {
		return false
	}
	return true
}

func (r *houseRoutes) create(c *gin.Context) {
	var request createHouseRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, r.l, http.StatusBadRequest, "http - createHouse - invalid request body", err)
		return
	}

	if !request.satisfiesTheConditions() {
		errorResponse(c, r.l, http.StatusBadRequest, "http - createHouse", errors.New("invalid parameter values"))
		return
	}

	house, err := r.h.Create(c.Request.Context(), entity.House{
		ID:        request.ID,
		Address:   request.Address,
		Year:      request.Year,
		Developer: request.Developer,
	})
	if err != nil {
		errorResponse(c, r.l, http.StatusInternalServerError, "http - createHouse - translation service problems", err)
		return
	}

	fmt.Println(house.CreatedAt)

	c.JSON(http.StatusOK, house)
}

type getFlatsResponse struct {
	Flats []entity.Flat `json:"flats"`
}

func (r *houseRoutes) getFlats(c *gin.Context) {
	houseID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, r.l, http.StatusBadRequest, "http - getFlats - invalid house ID", err)
		return
	}

	userType, exists := c.Get("userType")
	if !exists {
		errorResponse(c, r.l, http.StatusInternalServerError, "http - getFlats: user role not found", err)
		return
	}
	typeStr, ok := userType.(string)
	if !ok {
		errorResponse(c, r.l, http.StatusInternalServerError, "http - getFlats: convert error", err)
		return
	}

	flats, err := r.h.Flats(c.Request.Context(), houseID, typeStr)
	if err != nil {
		errorResponse(c, r.l, http.StatusInternalServerError, "http - getFlats: could not get flats", err)
		return
	}

	c.JSON(http.StatusOK, getFlatsResponse{Flats: flats})
}

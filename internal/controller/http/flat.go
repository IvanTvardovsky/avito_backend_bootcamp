package http

import (
	"avito_bootcamp/internal/entity"
	"avito_bootcamp/internal/usecase"
	"avito_bootcamp/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type flatRoutes struct {
	f usecase.Flat
	l logger.Interface
}

func newFlatRoutes(handler *gin.Engine, l logger.Interface, f usecase.Flat) {
	r := &flatRoutes{
		f: f,
		l: l,
	}

	h := handler.Group("/flat/")
	{
		h.POST("/create", RoleMiddleware([]string{"moderator", "client"}), r.create)
		h.POST("/update", RoleMiddleware([]string{"moderator"}), r.update)
	}
}

type createFlatRequest struct {
	Number  int `json:"number"`
	HouseID int `json:"house_id"`
	Price   int `json:"price"`
	Rooms   int `json:"rooms"`
}

func (f *createFlatRequest) satisfiesTheConditions() bool {
	if f.HouseID < 1 || f.Price < 0 || f.Rooms < 1 {
		fmt.Println(f)
		return false
	}
	return true
}

func (r *flatRoutes) create(c *gin.Context) {
	var request createFlatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, r.l, http.StatusBadRequest, "http - createFlat - invalid request body", err)
		return
	}

	flat, err := r.f.Create(c.Request.Context(), entity.Flat{
		Number:  request.Number,
		HouseID: request.HouseID,
		Price:   request.Price,
		Rooms:   request.Rooms,
	})
	if err != nil {
		errorResponse(c, r.l, http.StatusInternalServerError, "http - createFlat - translation service problems", err)
		return
	}

	c.JSON(http.StatusOK, flat)
}

type updateFlatRequest struct {
	ID      int    `json:"id"`
	Number  *int   `json:"number"`   // если поле не будет представлено в запросе, то указатель останется nil
	HouseID *int   `json:"house_id"` // что удобнее, чем работать с дефолтным 0 для int
	Price   *int   `json:"price"`
	Rooms   *int   `json:"rooms"`
	Status  string `json:"status"`
}

func (r *flatRoutes) update(c *gin.Context) {
	var request updateFlatRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		errorResponse(c, r.l, http.StatusBadRequest, "http - updateFlat - invalid request body", err)
		return
	}

	flatToUpdate := entity.Flat{
		ID: request.ID,
	}

	if request.Number != nil {
		flatToUpdate.Number = *request.Number
	} else {
		flatToUpdate.Number = -1
	}

	if request.HouseID != nil {
		flatToUpdate.HouseID = *request.HouseID
	} else {
		flatToUpdate.HouseID = -1
	}

	if request.Price != nil {
		flatToUpdate.Price = *request.Price
	} else {
		flatToUpdate.Price = -1
	}

	if request.Rooms != nil {
		flatToUpdate.Rooms = *request.Rooms
	} else {
		flatToUpdate.Rooms = -1
	}

	flatToUpdate.Status = request.Status

	var flat entity.Flat
	var err error

	if flat, err = r.f.Update(c.Request.Context(), flatToUpdate); err != nil {
		errorResponse(c, r.l, http.StatusInternalServerError, "http - updateFlat - could not update flat", err)
		return
	}

	c.JSON(http.StatusOK, flat)
}

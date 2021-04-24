package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raismaulana/ticketing-event/app/dto"
	"github.com/raismaulana/ticketing-event/app/entity"
	"github.com/raismaulana/ticketing-event/app/helper"
	"github.com/raismaulana/ticketing-event/app/usecase"
)

type EventController interface {
	Insert(c *gin.Context)
	Fetch(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Test(c *gin.Context)
}

type eventController struct {
	eventCase usecase.EventCase
	redisCase usecase.RedisCase
}

func NewEventController(eventCase usecase.EventCase, redisCase usecase.RedisCase) EventController {
	return &eventController{
		eventCase: eventCase,
		redisCase: redisCase,
	}
}

func (ctrl *eventController) Insert(c *gin.Context) {
	var insertEventDTO dto.InsertEventDTO
	if err := c.ShouldBind(&insertEventDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	event, err := ctrl.eventCase.Insert(insertEventDTO)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusCreated, helper.BuildResponse(true, "OK!", event))
}

func (ctrl *eventController) Fetch(c *gin.Context) {
	events, errRes := ctrl.eventCase.Fetch()

	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}

	ctrl.redisCase.Set(c.FullPath(), events)
	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", events))
}
func (ctrl *eventController) Test(c *gin.Context) {
	var events []entity.Event
	obj, err := ctrl.redisCase.Get("/api/event/", events)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", obj))
}
func (ctrl *eventController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	event, errRes := ctrl.eventCase.GetByID(id)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", event))
}

func (ctrl *eventController) Update(c *gin.Context) {
	var updateEventDTO dto.UpdateEventDTO
	if err := c.ShouldBind(&updateEventDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	event, errRes := ctrl.eventCase.Update(updateEventDTO)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusConflict, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", event))
}

func (ctrl *eventController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}
	_, errRes := ctrl.eventCase.Delete(id, time.Now())
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(true, "Deleted!", helper.EmptyObj{}))
}

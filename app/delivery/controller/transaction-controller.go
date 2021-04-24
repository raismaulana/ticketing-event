package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raismaulana/ticketing-event/app/dto"
	"github.com/raismaulana/ticketing-event/app/helper"
	"github.com/raismaulana/ticketing-event/app/usecase"
)

type TransactionController interface {
	Insert(c *gin.Context)
	Fetch(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type transactionController struct {
	transactionCase usecase.TransactionCase
	redisCase       usecase.RedisCase
}

func NewTransactionController(transactionCase usecase.TransactionCase, redisCase usecase.RedisCase) TransactionController {
	return &transactionController{
		transactionCase: transactionCase,
		redisCase:       redisCase,
	}
}

func (ctrl *transactionController) Insert(c *gin.Context) {
	var insertTransactionDTO dto.InsertTransactionDTO
	if err := c.ShouldBind(&insertTransactionDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	transaction, err := ctrl.transactionCase.Insert(insertTransactionDTO)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusCreated, helper.BuildResponse(true, "OK!", transaction))
}

func (ctrl *transactionController) Fetch(c *gin.Context) {
	transactions, errRes := ctrl.transactionCase.Fetch()

	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", transactions))
}

func (ctrl *transactionController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	transaction, errRes := ctrl.transactionCase.GetByID(id)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", transaction))
}

func (ctrl *transactionController) Update(c *gin.Context) {
	var updateTransactionDTO dto.UpdateTransactionDTO
	if err := c.ShouldBind(&updateTransactionDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	transaction, errRes := ctrl.transactionCase.Update(updateTransactionDTO)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusConflict, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", transaction))
}

func (ctrl *transactionController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}
	_, errRes := ctrl.transactionCase.Delete(id, time.Now())
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(true, "Deleted!", helper.EmptyObj{}))
}

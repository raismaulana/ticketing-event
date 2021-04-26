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

type UserController interface {
	Fetch(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	AllEventReport(c *gin.Context)
}

type userController struct {
	userCase  usecase.UserCase
	redisCase usecase.RedisCase
}

func NewUserController(userCase usecase.UserCase, redisCase usecase.RedisCase) UserController {
	return &userController{
		userCase:  userCase,
		redisCase: redisCase,
	}
}

func (ctrl *userController) Fetch(c *gin.Context) {
	users, errRes := ctrl.userCase.Fetch()

	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", users))
}

func (ctrl *userController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	user, errRes := ctrl.userCase.GetByID(id)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", user))
}

func (ctrl *userController) Update(c *gin.Context) {
	var updateUserDTO dto.UpdateUserDTO
	if err := c.ShouldBind(&updateUserDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	user, errRes := ctrl.userCase.Update(updateUserDTO)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusConflict, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", user))
}

func (ctrl *userController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}
	_, errRes := ctrl.userCase.Delete(id, time.Now())
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusOK, helper.BuildResponse(true, "Deleted!", helper.EmptyObj{}))
}

func (ctrl *userController) AllEventReport(c *gin.Context) {
	users, err := ctrl.userCase.AllEventReport()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}
	c.JSON(http.StatusOK, helper.BuildResponse(true, "ok", users))
}

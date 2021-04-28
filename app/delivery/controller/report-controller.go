package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raismaulana/ticketing-event/app/helper"
	"github.com/raismaulana/ticketing-event/app/usecase"
)

type ReportController interface {
	FetchAllReportUserBoughtEvent(c *gin.Context)
	FetchAllReportEventByCreator(c *gin.Context)
}

type reportController struct {
	reportCase usecase.ReportCase
}

func NewReportController(reportCase usecase.ReportCase) ReportController {
	return &reportController{
		reportCase: reportCase,
	}
}

func (ctrl *reportController) FetchAllReportUserBoughtEvent(c *gin.Context) {

	res, err := ctrl.reportCase.FetchAllReportUserBoughtEvent()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}
	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", helper.BuildResponse(true, "OK!", res)))
}

func (ctrl *reportController) FetchAllReportEventByCreator(c *gin.Context) {
	sortf := c.DefaultQuery("sortf", "t.id")
	sorta := c.DefaultQuery("sorta", "ASC")

	if !(sortf == "t.id" || sortf == "total_participant" || sortf == "total_amount") {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", "Bad Params", helper.EmptyObj{}))
		return
	}

	creatorId, ok := c.MustGet("user_id").(string)
	role, ok2 := c.MustGet("user_role").(string)
	if !(ok && ok2) {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", "Token invalid", helper.EmptyObj{}))
		return
	}
	// id, err := strconv.ParseUint(creatorId, 10, 64)
	// if err != nil {
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
	// 	return
	// }
	if role == "admin" {
		creatorId = "e.creator_id"
	}

	reportEvent, errRes := ctrl.reportCase.FetchAllReportEventByCreator(creatorId, sortf, sorta)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}
	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", reportEvent))
}

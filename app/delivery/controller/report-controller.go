package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raismaulana/ticketing-event/app/helper"
	"github.com/raismaulana/ticketing-event/app/usecase"
)

type ReportController interface {
	FetchAllReportEvent(c *gin.Context)
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

func (ctrl *reportController) FetchAllReportEvent(c *gin.Context) {

}

func (ctrl *reportController) FetchAllReportEventByCreator(c *gin.Context) {
	sortf := c.DefaultQuery("sortf", "t.id")
	sorta := c.DefaultQuery("sorta", "ASC")
	log.Println(sortf, sorta)
	creatorId, ok := c.MustGet("user_id").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", "Token invalid", helper.EmptyObj{}))
		return
	}
	id, err := strconv.ParseUint(creatorId, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}
	reportEvent, errRes := ctrl.reportCase.FetchAllReportEventByCreator(id, sortf, sorta)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusNoContent, helper.BuildErrorResponse("error", errRes.Error(), helper.EmptyObj{}))
		return
	}
	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", reportEvent))
}

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
	BuyEvent(c *gin.Context)
	UploadReceipt(c *gin.Context)
	VerifyPayment(c *gin.Context)
}

type transactionController struct {
	transactionCase usecase.TransactionCase
	userCase        usecase.UserCase
	redisCase       usecase.RedisCase
}

func NewTransactionController(transactionCase usecase.TransactionCase, userCase usecase.UserCase, redisCase usecase.RedisCase) TransactionController {
	return &transactionController{
		transactionCase: transactionCase,
		userCase:        userCase,
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

func (ctrl *transactionController) BuyEvent(c *gin.Context) {
	var buyEventDTO dto.BuyEventDTO

	if err := c.ShouldBind(&buyEventDTO); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	participantId, ok := c.MustGet("user_id").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", "Token invalid", helper.EmptyObj{}))
		return
	}
	parsedparticipantId, errParse := strconv.ParseUint(participantId, 10, 64)
	if errParse != nil {
		panic(errParse)
	}

	buyEventDTO.ParticipantId = parsedparticipantId
	transaction, errRes := ctrl.transactionCase.BuyEvent(buyEventDTO)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusConflict, helper.BuildErrorResponse("Failed Buy Event", errRes.Error(), transaction))
		return
	}
	user, err := ctrl.userCase.GetByID(parsedparticipantId)
	if err == nil {
		go helper.SendMail(user.Email, "Pay Your Event ticket", "You can transafer using this method")
	}
	c.JSON(http.StatusCreated, helper.BuildResponse(true, "Check out success, check your email", transaction))
}

func (ctrl *transactionController) UploadReceipt(c *gin.Context) {
	participantId, ok := c.MustGet("user_id").(string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", "Token invalid", helper.EmptyObj{}))
		return
	}

	var uploadReceipt dto.UploadReceipt

	if err := c.ShouldBind(&uploadReceipt); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	parsedparticipantId, errParse := strconv.ParseUint(participantId, 10, 64)
	if errParse != nil {
		panic(errParse)
	}
	uploadReceipt.ParticipantId = parsedparticipantId

	transaction, errRes := ctrl.transactionCase.UploadReceipt(uploadReceipt)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.BuildErrorResponse("Failed Upload", errRes.Error(), transaction))
		return
	}

	c.JSON(http.StatusCreated, helper.BuildResponse(true, "Upload Success, Wait for verified", transaction))
}

func (ctrl *transactionController) VerifyPayment(c *gin.Context) {
	var verify dto.Verify
	if err := c.ShouldBind(&verify); err != nil || verify.Status != "passed" && verify.Status != "failed" {
		c.AbortWithStatusJSON(http.StatusBadRequest, helper.BuildErrorResponse("error", err.Error(), helper.EmptyObj{}))
		return
	}

	res, err2 := ctrl.transactionCase.VerifyPayment(verify)
	if err2 != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, helper.BuildErrorResponse("error", err2.Error(), helper.EmptyObj{}))
		return
	}
	c.JSON(http.StatusOK, helper.BuildResponse(true, "OK!", helper.BuildResponse(true, "Success Update Data", res)))

}

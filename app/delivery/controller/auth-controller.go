package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raismaulana/ticketing-event/app/dto"
	"github.com/raismaulana/ticketing-event/app/entity"
	"github.com/raismaulana/ticketing-event/app/helper"
	"github.com/raismaulana/ticketing-event/app/usecase"
)

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type authController struct {
	authCase usecase.AuthCase
	jwtCase  usecase.JWTCase
}

func NewAuthController(authCase usecase.AuthCase, jwtCase usecase.JWTCase) AuthController {
	return &authController{
		authCase: authCase,
		jwtCase:  jwtCase,
	}
}

func (ctrl *authController) Login(c *gin.Context) {
	var loginDTO dto.LoginDTO
	err := c.ShouldBind(&loginDTO)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to proccess request", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	result := ctrl.authCase.Login(loginDTO)
	if v, ok := result.(entity.User); ok {
		v.Token = ctrl.jwtCase.GenerateToken(strconv.FormatUint(v.ID, 10), v.Role)
		log.Println(v)
		response := helper.BuildResponse(true, "Login Success!", v)
		c.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorResponse("Login Failed!", "Username or Password is wrong!", helper.EmptyObj{})
	c.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (ctrl *authController) Register(c *gin.Context) {
	var registerDTO dto.RegisterUserDTO
	if err := c.ShouldBind(&registerDTO); err != nil {
		response := helper.BuildErrorResponse("Failed to proccess request", err.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	_, err1 := ctrl.authCase.IsDuplicateUnique(registerDTO.Username, registerDTO.Email)
	if err1 != nil {
		response := helper.BuildErrorResponse("Register Failed!", err1.Error(), helper.EmptyObj{})
		c.JSON(http.StatusConflict, response)
		return
	}

	user, err2 := ctrl.authCase.Register(registerDTO)
	if err2 != nil {
		response := helper.BuildErrorResponse("Register Failed!", err2.Error(), helper.EmptyObj{})
		c.AbortWithStatusJSON(http.StatusConflict, response)
		return
	}
	response := helper.BuildResponse(true, "Register Success!", user)
	c.JSON(http.StatusCreated, response)
}

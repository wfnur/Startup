package handler

import (
	"net/http"
	"startup/helper"
	"startup/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ErrorValidationFormatter(err)
		errMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register Failed", http.StatusUnprocessableEntity, "Error", errMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	registerUser, err := h.userService.RegisterUser(input)
	if err != nil {
		errMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Register Failed", http.StatusBadRequest, "Error", errMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(registerUser, "tokentokentoken")

	response := helper.APIResponse("Account has been created", http.StatusOK, "Success", formatter)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	//validasi input
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ErrorValidationFormatter(err)
		errMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "Error", errMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	//input di proses
	loginUser, err := h.userService.Login(input)
	if err != nil {
		errMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login Failed", http.StatusBadRequest, "Error", errMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//return input
	formatter := user.FormatUser(loginUser, "tokentokentoken")

	response := helper.APIResponse("Login Successfully", http.StatusOK, "Success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckAvailabilityEmail(c *gin.Context) {
	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ErrorValidationFormatter(err)
		errMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Checking Email Failed", http.StatusUnprocessableEntity, "Error", errMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	IsEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Check Email Failed", http.StatusBadRequest, "Error", errMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//return input
	data := gin.H{
		"is_available": IsEmailAvailable,
	}

	var message string
	if IsEmailAvailable {
		message = "Email Is Available"
	} else {
		message = "Email Registered"
	}

	response := helper.APIResponse(message, http.StatusOK, "Success", data)
	c.JSON(http.StatusOK, response)

}

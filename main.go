package main

import (
	"log"
	"startup/handler"
	"startup/user"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:kangbaso@tcp(127.0.0.1:3306)/startup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	//userService.SaveAvatar(9, "images/1-img.jpg")

	router := gin.Default()
	api := router.Group("/api/v1")
	api.POST("/user", userHandler.RegisterUser)
	api.POST("/session", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckAvailabilityEmail)
	api.POST("/avatars", userHandler.UploadAvatar)

	router.Run()
	// userInput := user.RegisterUserInput{}
	// userInput.Name = "Test Simpan dari service"
	// userInput.Email = "service@gmail.com"
	// userInput.Occupation = "Service"
	// userInput.Password = "password"
	// userService.RegisterUser(userInput)

}

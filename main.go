package main

import (
	"fmt"
	"log"
	"net/http"
	"startup/auth"
	"startup/campaign"
	"startup/handler"
	"startup/helper"
	"startup/user"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	campaignRepo := campaign.NewRepository(db)

	/*
		======== TEST CAMPAIGN =======
		-- test findall dan findby userid
	*/
	campaigns, _ := campaignRepo.FindAll()

	for _, campaign := range campaigns {
		fmt.Println(campaign.Name)
		if len(campaign.CampaignImages) > 0 {
			fmt.Println(campaign.CampaignImages[0].FileName)
		}
	}
	fmt.Println("====")
	campaignsByUserid, _ := campaignRepo.FindByUserID(11)
	for _, campaign := range campaignsByUserid {
		fmt.Println(campaign.Name)
		if len(campaign.CampaignImages) > 0 {
			fmt.Println(campaign.CampaignImages[0].FileName)
		}

	}

	userService := user.NewService(userRepo)
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")
	api.POST("/user", userHandler.RegisterUser)
	api.POST("/session", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckAvailabilityEmail)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	router.Run()
	// userInput := user.RegisterUserInput{}
	// userInput.Name = "Test Simpan dari service"
	// userInput.Email = "service@gmail.com"
	// userInput.Occupation = "Service"
	// userInput.Password = "password"
	// userService.RegisterUser(userInput)

}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		//get Header
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		user_id := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(user_id)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("CurrentUser", user)

	}
}

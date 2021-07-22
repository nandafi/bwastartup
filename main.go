package main

import (
	"bwastartup/auth"
	"bwastartup/campaign"
	"bwastartup/handler"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "debian-sys-maint:uCSlljDvEyR7YJ0H@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	//Inisialisasi repositori
	userRepository := user.NewRepository(db)

	campaignRepository := campaign.NewRepository(db)
	campaigns, err := campaignRepository.FindByUserID(7)
	fmt.Println("debug")
	fmt.Println("debug")
	fmt.Println("debug")
	fmt.Println(len(campaigns))
	for _, campaign := range campaigns {
		fmt.Println(campaign.Name)
		if len(campaign.CampaignImages) > 0 {
			fmt.Println("jumlah gambar")
			fmt.Println(len(campaign.CampaignImages))
			fmt.Println(campaign.CampaignImages[0].FileName)
		}
	}

	userService := user.NewService(userRepository)
	authService := auth.NewService()

	//register
	userHandler := handler.NewUserHandler(userService, authService)
	router := gin.Default()
	api := router.Group("api/v1")
	api.POST("/users", userHandler.RegisterUser)                                             //register
	api.POST("/sessions", userHandler.Login)                                                 //login
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)                          //check email available apa nggk
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar) //upload avatar

	router.Run()

}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {

	//saat panggil func authMiddleware ini, ntar nilai kembalian adalah sebuah func dibawah ini
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		//jika didalam authHeader tidak ada string "Bearer" maka....
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		//formatnya = Bearer(spasi)tokenxxxxx
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ") //split pakai spasi
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

		//jika (tidak oke) atau (token tidak valid) maka...
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		c.Set("currentUser", user)

	}

}

//---middleware
//ambil nilai header Authorization: Bearer tokentokentoken
//dari header Authorization, kita ambil nilai tokennya saja
//kita validasi token
//kita ambil user_id
//ambil user dari db berdasarkan user_id lewat service
//kita set context isinya user

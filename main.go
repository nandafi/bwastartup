package main

import (
	"bwastartup/handler"
	"bwastartup/user"
	"log"

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
	userService := user.NewService(userRepository)
	//register
	userHandler := handler.NewUserHandler(userService)
	router := gin.Default()
	api := router.Group("api/v1")
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	router.Run()

}

//input
//handler
//service
//repository
//db

package main

import (
	"DemoAppBE/controllers"
	"DemoAppBE/models"

	"github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    db := models.SetupModels()

    // Provide db variable to controllers
    router.Use(func(c *gin.Context) {
        c.Set("db", db)
        c.Next()
    })

    router.GET("/home", controllers.HomeItems)
    router.GET("/item/:id", controllers.FindItem) // new
    router.GET("/items", controllers.FindItems) // new
    router.GET("/items/search", controllers.SearchItems)

    router.GET("/favorites/:user_id", controllers.ReadFavorites)
    router.POST("/favorites", controllers.InsertFavorites)

    router.Run()
}
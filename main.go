package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shubhamxg/go-hunger/controller"
)

// func init() {
// //TODO: env variables are hardcoded in code i will write code to import them from env file
// }

func main() {
	router := gin.Default()
	handler := controller.NewRecipesHandler()
	router.GET("/recipes/", handler.ListRecipesHandler)
	router.GET("/recipes/:id", handler.GetRecipeHandler)
	router.POST("/recipes", handler.NewRecipeHandler)
	router.PUT("/recipes/:id", handler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", handler.DeleteRecipeHandler)
	router.GET("/recipes/search", handler.SearchRecipeHandler)
	router.Run()
}

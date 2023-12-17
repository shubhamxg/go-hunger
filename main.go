package main

import (
	// "context"
	// "fmt"

	"github.com/gin-gonic/gin"
	// "github.com/redis/go-redis/v9"
	"github.com/shubhamxg/go-hunger/controller"
)

func init() {
	// TODO: env variables are hardcoded in code i will write code to import them from env file
	// redis_client := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "",
	// 	DB:       0,
	// })
	// status := redis_client.Ping(context.TODO())
	// fmt.Println(status)
}

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

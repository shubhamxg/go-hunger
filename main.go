package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/sessions"
	// "github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	redis_store "github.com/gin-contrib/sessions/redis"
	"github.com/shubhamxg/go-hunger/controller"
	"github.com/shubhamxg/go-hunger/models"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Failed to load .env file")
	}
}

func main() {
	router := gin.Default()
	handler := controller.NewRecipesHandler()
	auth_handler := controller.NewAuthHandler()
	store, _ := redis_store.NewStore(
		10,
		"tcp",
		fmt.Sprintf(
			"%s:%s",
			models.Env(models.REDIS_HOST),
			models.Env(models.REDIS_PORT),
		),
		models.Env(models.REDIS_PASSWORD),
		[]byte("secret"),
	)
	router.Use(sessions.Sessions("recipe_api", store))

	router.GET("/recipes/", handler.ListRecipesHandler)
	router.GET("/recipes/:id", handler.GetRecipeHandler)
	router.GET("/recipes/search", handler.SearchRecipeHandler)
	router.POST("/signin", auth_handler.SignInHandler)
	router.POST("/signup", auth_handler.SignUpHandler)
	router.POST("/signout", auth_handler.SignOutHandler)
	router.POST("/refresh", auth_handler.RefreshHandler)

	authorized := router.Group("/")
	authorized.Use(auth_handler.AuthMiddlerware())
	{
		authorized.POST("/recipes", handler.NewRecipeHandler)
		authorized.PUT("/recipes/:id", handler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", handler.DeleteRecipeHandler)
	}

	// run_addr := "localhost"
	// run_port := ":8080"
	// log.Printf("Running on %s%s", run_addr, run_port)
	// if err := router.Run(run_port); err != nil {
	// 	fmt.Printf("Error > Gin Failed to run : %v", err.Error())
	// 	return
	// }
	router.Run()
}

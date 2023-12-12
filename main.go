package main

import (
	"encoding/json"
	"fmt"
	// "net/http"
	"os"
	// "time"

	"github.com/gin-gonic/gin"
	// "github.com/rs/xid"
)

type Recipe struct {
	Recipe RecipeData `json:"recipe"`
}

type RecipeData struct {
	Name             string            `json:"name"`
	ID               string            `json:"id"`
	Description      string            `json:"description"`
	Tags             []string          `json:"tag"`
	Ingredients      []Ingredient      `json:"ingredient"`
	IngredientGroups []IngredientGroup `json:"ingredientGroup"`
	Steps            []Step            `json:"step"`
}

type Ingredient struct {
	Name        string `json:"name"`
	Amount      string `json:"amount,omitempty"`
	Unit        string `json:"unit,omitempty"`
	Preparation string `json:"preparation,omitempty"`
}

type IngredientGroup struct {
	Name        string       `json:"name"`
	Ingredients []Ingredient `json:"ingredient"`
}

type Step struct {
	Description string `json:"description"`
}

// var recipes []Recipe

func init() {
	// recipes = make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	// fmt.Println(string(file))
	if err := json.Unmarshal([]byte(file), &Recipe); err != nil {
		fmt.Println("Error >", err)
	}
}

// func NewRecipeHandler(c *gin.Context) {
// 	var recipe Recipe
// 	if err := c.ShouldBindJSON(&recipe); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	recipe.ID = xid.New().String()
// 	recipe.PublishedAt = time.Now()
// 	recipes = append(recipes, recipe)
// 	c.JSON(http.StatusOK, recipe)
// }

// func ListRecipeHandler(c *gin.Context) {
// 	c.JSON(http.StatusOK, recipes)
// }

func main() {
	router := gin.Default()
	// router.POST("/recipes", NewRecipeHandler)
	// router.GET("/recipes", ListRecipeHandler)
	router.Run()
}

package main

import (
	// "fmt"
	"encoding/json"
	"net/http"
	"time"

	// "strings"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/shubhamxg/go-hunger/models"
)

type Recipe struct {
	Id          int    `db:"id"`
	Recipe_data []byte `db:"recipe_data"`
}
type recipe struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Tags         []string `json:"tags"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Publishedat  string   `json:"publishedat"`
}

func NewRecipeHandler(c *gin.Context) {
	new_recipe := recipe{}
	if err := c.ShouldBindJSON(&new_recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	new_recipe.Id = xid.New().String()
	new_recipe.Publishedat = time.Now().String()

	db := models.Start()
	tx := db.MustBegin()
	foo, err := json.Marshal(new_recipe)
	if err != nil {
		panic(err)
	}
	tx.MustExec(
		`INSERT INTO recipes (recipe_data) VALUES ($1::jsonb);`,
		string(foo),
	)
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe Added Successfully",
	})
}

func UpdateRecipeHandler(c *gin.Context) {
	recipe_id := c.Param("id")
	updated_recipe := recipe{}
	if err := c.ShouldBindJSON(&updated_recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	updated_recipe_json, err := json.Marshal(updated_recipe)
	if err != nil {
		panic(err)
	}

	db := models.Start()
	tx := db.MustBegin()
	foo, err := db.Exec(
		`UPDATE Recipes SET recipe_data = $1::jsonb WHERE recipe_data ->> 'id' = $2;`,
		updated_recipe_json,
		recipe_id,
	)

	if err == nil {
		count, _ := foo.RowsAffected()
		if count == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Recipe not found",
			})
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe Updated Successfully",
	})
}

func DeleteRecipeHandler(c *gin.Context) {
	recipe_id := c.Param("id")
	db := models.Start()
	tx := db.MustBegin()
	foo, err := db.Exec(`DELETE FROM recipes WHERE recipe_data ->> 'id' = $1;`, recipe_id)
	if err == nil {
		count, _ := foo.RowsAffected()
		if count == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Recipe not found",
			})
			return
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe deleted Successfully",
	})
}

func SearchRecipeHandler(c *gin.Context) {
	search_tag := c.Query("tag")
	db := models.Start()

	recipes := []Recipe{}
	if err := db.Select(&recipes, `SELECT * FROM recipes 
WHERE EXISTS (
    SELECT * 
    FROM jsonb_array_elements_text(recipe_data->'tags') as tag
    WHERE tag = $1
);
`, search_tag); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	filtered_recipes := make([]recipe, 0)
	if len(recipes) > 0 {
		for i := 0; i < len(recipes); i++ {
			single_recipe := recipe{}
			_ = json.Unmarshal(recipes[i].Recipe_data, &single_recipe)
			filtered_recipes = append(filtered_recipes, single_recipe)

		}
		c.JSON(http.StatusOK, filtered_recipes)
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Recipe not found",
		})
		return
	}
}

func GetRecipeHandler(c *gin.Context) {
	recipe_id := c.Param("id")
	db := models.Start()

	recipes := []Recipe{}
	if err := db.Select(&recipes, `SELECT * FROM Recipes
WHERE recipe_data ->> 'id' = $1;
`, recipe_id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(recipes) > 0 {
		single_recipe := recipe{}
		_ = json.Unmarshal(recipes[0].Recipe_data, &single_recipe)
		c.JSON(http.StatusOK, single_recipe)
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Recipe not found",
		})
		return
	}
}

func ListRecipesHandler(c *gin.Context) {
	db := models.Start()
	recipes := []Recipe{}
	if err := db.Select(&recipes, `SELECT * FROM Recipes`); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	filtered_recipes := make([]recipe, 0)
	if len(recipes) > 0 {
		for i := 0; i < len(recipes); i++ {
			single_recipe := recipe{}
			_ = json.Unmarshal(recipes[i].Recipe_data, &single_recipe)
			filtered_recipes = append(filtered_recipes, single_recipe)
		}
		c.JSON(http.StatusOK, filtered_recipes)
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Recipe not found",
		})
		return
	}
}

func main() {
	router := gin.Default()
	router.GET("/recipes/", ListRecipesHandler)
	router.GET("/recipes/:id", GetRecipeHandler)
	router.POST("/recipes", NewRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)
	router.Run()
}

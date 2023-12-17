package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/shubhamxg/go-hunger/models"
)

type Inc int

const (
	not_found Inc = iota
	backend_error
	updated
	deleted
	added
)

func recipe_response(resp Inc) gin.H {
	switch resp {
	case not_found:
		return gin.H{"message": "Recipes not found"}
	case backend_error:
		return gin.H{"error": "Something went wrong at the backend"}
	case updated:
		return gin.H{"message": "Recipe Updated Successfully"}
	case added:
		return gin.H{"message": "Recipe Added Successfully"}
	case deleted:
		return gin.H{"message": "Recipe deleted Successfully"}
	default:
		return gin.H{"error": "Recipe Response not found"}
	}
}

type RecipeHandler struct {
	db *sqlx.DB
}

func NewRecipesHandler() *RecipeHandler {
	return &RecipeHandler{
		db: models.Start(),
	}
}

func (handler *RecipeHandler) NewRecipeHandler(c *gin.Context) {
	new_recipe := models.Recipe{}
	if err := c.ShouldBindJSON(&new_recipe); err != nil {
		c.JSON(http.StatusBadRequest, recipe_response(backend_error))
	}
	new_recipe.Id = xid.New().String()
	new_recipe.Publishedat = time.Now().String()

	tx := handler.db.MustBegin()
	json_new_recipe, err := json.Marshal(new_recipe)
	if err != nil {
		panic(err)
	}
	tx.MustExec(`INSERT INTO recipes (recipe_data) VALUES ($1::jsonb);`, string(json_new_recipe))
	tx.Commit()

	c.JSON(http.StatusOK, recipe_response(added))
}

func (handler *RecipeHandler) UpdateRecipeHandler(c *gin.Context) {
	recipe_id := c.Param("id")
	updated_recipe := models.Recipe{}
	if err := c.ShouldBindJSON(&updated_recipe); err != nil {
		c.JSON(http.StatusBadRequest, recipe_response(backend_error))
		return
	}

	updated_recipe_json, err := json.Marshal(updated_recipe)
	if err != nil {
		panic(err)
	}

	tx := handler.db.MustBegin()
	executed, err := handler.db.Exec(
		`UPDATE Recipes SET recipe_data = $1::jsonb WHERE recipe_data ->> 'id' = $2;`,
		updated_recipe_json,
		recipe_id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": recipe_response(backend_error),
		})
		return
	}

	if count, _ := executed.RowsAffected(); count == 0 {
		c.JSON(http.StatusNotFound, recipe_response(not_found))
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, recipe_response(updated))
}

func (handler *RecipeHandler) DeleteRecipeHandler(c *gin.Context) {
	recipe_id := c.Param("id")
	tx := handler.db.MustBegin()
	executed, err := handler.db.Exec(
		`DELETE FROM recipes WHERE recipe_data ->> 'id' = $1;`,
		recipe_id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, recipe_response(backend_error))
	}
	if count, _ := executed.RowsAffected(); count == 0 {
		c.JSON(http.StatusNotFound, recipe_response(not_found))
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, recipe_response(deleted))
}

func (handler *RecipeHandler) SearchRecipeHandler(c *gin.Context) {
	search_tag := c.Query("tag")

	recipes := []models.Recipes{}
	if err := handler.db.Select(&recipes, `
		SELECT * FROM recipes 
		WHERE EXISTS (SELECT * FROM jsonb_array_elements_text(recipe_data->'tags') as tag WHERE tag = $1);`, search_tag); err != nil {
		c.JSON(http.StatusNotFound, recipe_response(backend_error))
		return
	}

	filtered_recipes := make([]models.Recipe, 0)
	if len(recipes) > 0 {
		for i := 0; i < len(recipes); i++ {
			single_recipe := models.Recipe{}
			_ = json.Unmarshal(recipes[i].Recipe_data, &single_recipe)
			filtered_recipes = append(filtered_recipes, single_recipe)
		}
		c.JSON(http.StatusOK, filtered_recipes)
		return
	}
	c.JSON(http.StatusNotFound, recipe_response(not_found))
}

func (handler *RecipeHandler) GetRecipeHandler(c *gin.Context) {
	recipe_id := c.Param("id")

	recipes := []models.Recipes{}
	get_recipe_query := fmt.Sprintf(
		`SELECT * FROM Recipes WHERE recipe_data ->> 'id' = '%s';`,
		recipe_id,
	)
	if err := handler.db.Select(&recipes, get_recipe_query); err != nil {
		c.JSON(http.StatusInternalServerError, recipe_response(backend_error))
		return
	}

	if len(recipes) > 0 {
		single_recipe := models.Recipe{}
		_ = json.Unmarshal(recipes[0].Recipe_data, &single_recipe)
		c.JSON(http.StatusOK, single_recipe)
		return
	}
	c.JSON(http.StatusNotFound, recipe_response(not_found))
}

func (handler *RecipeHandler) ListRecipesHandler(c *gin.Context) {
	recipes := []models.Recipes{}
	if err := handler.db.Select(&recipes, `SELECT * FROM Recipes`); err != nil {
		c.JSON(http.StatusInternalServerError, recipe_response(backend_error))
		return
	}

	filtered_recipes := make([]models.Recipe, 0)
	if len(recipes) > 0 {
		for i := 0; i < len(recipes); i++ {
			single_recipe := models.Recipe{}
			_ = json.Unmarshal(recipes[i].Recipe_data, &single_recipe)
			filtered_recipes = append(filtered_recipes, single_recipe)
		}
		c.JSON(http.StatusOK, filtered_recipes)
		return
	}
	c.JSON(http.StatusNotFound, recipe_response(not_found))
}

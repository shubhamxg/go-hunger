package main

import (
	// "encoding/json"
	"fmt"
	// "os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shubhamxg/go-hunger/models"
)

var recipes []recipe

type recipe struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	Publishedat  time.Time `json:"publishedat"`
}

var schema = `
	CREATE TABLE recipes (
	id SERIAL PRIMARY KEY,
	recipe_data JSONB NOT NULL
	);
`

var schema_user = `
	CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL
	);
`

func main() {
	// recipes = make([]recipe, 0)
	// file, _ := os.ReadFile("recipes.json")
	// _ = json.Unmarshal([]byte(file), &recipes)
	db := models.Start()
	db.MustExec(schema_user)
	// for i := 0; i < len(recipes); i++ {
	// 	Recipe(recipes[i], db)
	// }
	// GetRecipe(db)
}

// func Recipe(recipe recipe, db *sqlx.DB) {
// 	// fmt.Println(recipe)
//
// 	tx := db.MustBegin()
//
// 	foo, err := json.Marshal(recipe)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// fmt.Println(string(foo))
// 	tx.MustExec(
// 		`INSERT INTO recipes (recipe_data) VALUES ($1::jsonb);`,
// 		string(foo),
// 	)
// 	tx.Commit()
// }

func GetRecipe(db *sqlx.DB) {
	type foo struct {
		Id          int    `db:"id"`
		Recipe_data []byte `db:"recipe_data"`
	}
	rec := []foo{}
	if err := db.Select(&rec, `SELECT * FROM recipes
WHERE recipe_data ->> 'name' = 'Oregano Marinated Chicken';
`); err != nil {
		fmt.Println("Something went wrong in getting recipe")
		panic(err)
	}
	fmt.Println(string(rec[0].Recipe_data))
}

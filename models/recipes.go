package models

type Recipes struct {
	Id          int    `db:"id"`
	Recipe_data []byte `db:"recipe_data"`
}

type Recipe struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Tags         []string `json:"tags"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Publishedat  string   `json:"publishedat"`
}

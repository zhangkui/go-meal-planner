package model

import "time"

type Ingredient struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
}

type Recipe struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Servings    int          `json:"servings"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
}

type MealType string

const (
	Breakfast MealType = "breakfast"
	Lunch     MealType = "lunch"
	Dinner    MealType = "dinner"
)

func (m MealType) Valid() bool {
	return m == Breakfast || m == Lunch || m == Dinner
}

type MenuEntry struct {
	Date      time.Time `json:"date"`
	Meal      MealType  `json:"meal"`
	RecipeID  string    `json:"recipe_id"`
	Servings  int       `json:"servings"`
	Confirmed bool      `json:"confirmed"`
}

type InventoryItem struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type ShoppingItem struct {
	Name      string  `json:"name"`
	Quantity  float64 `json:"quantity"`
	Unit      string  `json:"unit"`
	Purchased bool    `json:"purchased"`
}

type NutritionSummary struct {
	Calories        float64 `json:"calories"`
	Protein         float64 `json:"protein"`
	IngredientKinds int     `json:"ingredient_kinds"`
}

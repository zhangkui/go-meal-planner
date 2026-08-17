package service

import (
	"testing"
	"time"

	"github.com/zhangkui/go-meal-planner/internal/model"
)

func TestShoppingListSubtractsInventoryAfterAggregation(t *testing.T) {
	recipes, menus, inventory, shopping, _ := testServices()
	recipe := sampleRecipe("Egg meal")
	recipe.Ingredients = []model.Ingredient{{Name: "Egg", Quantity: 2, Unit: "piece", Calories: 140, Protein: 12}}
	created, _ := recipes.Create(recipe)
	_, _ = menus.Plan("2026-08-17", model.Breakfast, created.ID, 2, true)
	_, _ = menus.Plan("2026-08-18", model.Breakfast, created.ID, 2, true)
	_, _ = inventory.Add("Egg", "piece", 3)
	from, _ := time.Parse("2006-01-02", "2026-08-17")
	to, _ := time.Parse("2006-01-02", "2026-08-18")
	items, err := shopping.Generate(from, to)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Name != "Egg" || items[0].Quantity != 1 {
		t.Fatalf("shopping items = %+v, want one egg", items)
	}
}

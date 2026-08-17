package service

import (
	"errors"
	"testing"
	"time"

	"github.com/zhangkui/go-meal-planner/internal/model"
	"github.com/zhangkui/go-meal-planner/internal/store"
)

func testServices() (*RecipeService, *MenuService, *InventoryService, *ShoppingService, *NutritionService) {
	memory := store.NewMemory()
	recipes := NewRecipeService(memory)
	menus := NewMenuService(memory, recipes)
	inventory := NewInventoryService(memory)
	return recipes, menus, inventory, NewShoppingService(memory, recipes, menus, inventory), NewNutritionService(recipes, menus)
}

func sampleRecipe(name string) model.Recipe {
	return model.Recipe{
		Name:        name,
		Servings:    2,
		Ingredients: []model.Ingredient{{Name: "Rice", Quantity: 200, Unit: "g", Calories: 720, Protein: 14}},
		Steps:       []string{"Cook"},
	}
}

func TestRecipeCRUD(t *testing.T) {
	recipes, _, _, _, _ := testServices()
	created, err := recipes.Create(sampleRecipe("Rice bowl"))
	if err != nil {
		t.Fatal(err)
	}
	created.Name = "Changed outside"
	got, err := recipes.Get(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "Rice bowl" {
		t.Fatalf("stored recipe changed: %q", got.Name)
	}
	if err := recipes.Delete(created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := recipes.Get(created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestMenuAndNutritionForDate(t *testing.T) {
	recipes, menus, _, _, nutrition := testServices()
	recipe, err := recipes.Create(sampleRecipe("Lunch"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := menus.Plan("2026-08-17", model.Lunch, recipe.ID, 1, true); err != nil {
		t.Fatal(err)
	}
	date, _ := time.Parse("2006-01-02", "2026-08-17")
	summary, err := nutrition.ForDate(date)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Calories != 360 || summary.Protein != 7 || summary.IngredientKinds != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
}

func TestInventoryAddAdjustAndInsufficient(t *testing.T) {
	_, _, inventory, _, _ := testServices()
	if _, err := inventory.Add("Milk", "ml", 500); err != nil {
		t.Fatal(err)
	}
	if _, err := inventory.Consume("Milk", 200); err != nil {
		t.Fatal(err)
	}
	if _, err := inventory.Consume("Milk", 400); !errors.Is(err, ErrInsufficientStock) {
		t.Fatalf("expected insufficient error, got %v", err)
	}
	item, err := inventory.Adjust("Milk", "ml", 100)
	if err != nil {
		t.Fatal(err)
	}
	if item.Quantity != 100 {
		t.Fatalf("unexpected quantity: %v", item.Quantity)
	}
}

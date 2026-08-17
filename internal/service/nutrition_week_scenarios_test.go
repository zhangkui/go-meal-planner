package service

import (
	"testing"
	"time"

	"github.com/zhangkui/go-meal-planner/internal/model"
)

// TestWeeklyNutritionSundayBoundary verifies that querying the week on a
// Sunday only counts Monday through Sunday of that same week and never
// leaks the following Monday, even when meals exist on the boundary days.
func TestWeeklyNutritionSundayBoundary(t *testing.T) {
	recipes, menus, _, _, nutrition := testServices()
	recipe, _ := recipes.Create(sampleRecipe("Sunday boundary meal"))
	_, _ = menus.Plan("2026-08-17", model.Lunch, recipe.ID, 2, true)  // Monday this week
	_, _ = menus.Plan("2026-08-23", model.Dinner, recipe.ID, 2, true) // Sunday this week
	_, _ = menus.Plan("2026-08-24", model.Lunch, recipe.ID, 2, true)  // Monday next week
	date, _ := time.Parse("2006-01-02", "2026-08-23") // query on Sunday
	summary, err := nutrition.ForWeek(date)
	if err != nil {
		t.Fatal(err)
	}
	// Two meals this week (Mon + Sun); the following Monday must be excluded.
	if summary.Calories != 1440 || summary.Protein != 28 {
		t.Fatalf("weekly summary leaks next week when querying on Sunday: %+v", summary)
	}
}

// TestWeeklyNutritionExcludesDifferentIngredientOnNextMonday verifies that a
// meal on the following Monday with a different ingredient does not inflate
// calories, protein, or the ingredient-kind count for the current week.
func TestWeeklyNutritionExcludesDifferentIngredientOnNextMonday(t *testing.T) {
	recipes, menus, _, _, nutrition := testServices()
	rice, _ := recipes.Create(sampleRecipe("Rice meal"))
	beans, _ := recipes.Create(model.Recipe{
		Name:        "Beans meal",
		Servings:    2,
		Ingredients: []model.Ingredient{{Name: "Beans", Quantity: 150, Unit: "g", Calories: 500, Protein: 20}},
		Steps:       []string{"Cook"},
	})
	_, _ = menus.Plan("2026-08-17", model.Lunch, rice.ID, 2, true)  // Monday this week: Rice
	_, _ = menus.Plan("2026-08-24", model.Lunch, beans.ID, 2, true) // Monday next week: Beans
	date, _ := time.Parse("2006-01-02", "2026-08-19") // query mid-week
	summary, err := nutrition.ForWeek(date)
	if err != nil {
		t.Fatal(err)
	}
	// Only Rice this week; Beans on the following Monday must not contribute.
	if summary.Calories != 720 || summary.Protein != 14 || summary.IngredientKinds != 1 {
		t.Fatalf("weekly summary leaks next Monday ingredient: %+v", summary)
	}
}

// TestWeeklyNutritionOnlyCurrentWeek verifies the normal case where only the
// current week's menus exist: every day Monday through Sunday is counted.
func TestWeeklyNutritionOnlyCurrentWeek(t *testing.T) {
	recipes, menus, _, _, nutrition := testServices()
	recipe, _ := recipes.Create(sampleRecipe("Week meal"))
	_, _ = menus.Plan("2026-08-17", model.Breakfast, recipe.ID, 2, true) // Monday
	_, _ = menus.Plan("2026-08-19", model.Lunch, recipe.ID, 2, true)     // Wednesday
	_, _ = menus.Plan("2026-08-23", model.Dinner, recipe.ID, 2, true)    // Sunday
	date, _ := time.Parse("2006-01-02", "2026-08-19")
	summary, err := nutrition.ForWeek(date)
	if err != nil {
		t.Fatal(err)
	}
	// Three meals this week, all counted.
	if summary.Calories != 2160 || summary.Protein != 42 || summary.IngredientKinds != 1 {
		t.Fatalf("weekly summary miscounts current week: %+v", summary)
	}
}

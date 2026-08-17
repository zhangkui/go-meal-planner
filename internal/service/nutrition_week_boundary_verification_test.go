package service

import (
	"testing"
	"time"

	"github.com/zhangkui/go-meal-planner/internal/model"
)

func TestWeeklyNutritionExcludesFollowingMonday(t *testing.T) {
	recipes, menus, _, _, nutrition := testServices()
	recipe, _ := recipes.Create(sampleRecipe("Monday meal"))
	_, _ = menus.Plan("2026-08-17", model.Lunch, recipe.ID, 2, true)
	_, _ = menus.Plan("2026-08-24", model.Lunch, recipe.ID, 2, true)
	date, _ := time.Parse("2006-01-02", "2026-08-19")
	summary, err := nutrition.ForWeek(date)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Calories != 720 || summary.Protein != 14 {
		t.Fatalf("weekly summary includes another week: %+v", summary)
	}
}

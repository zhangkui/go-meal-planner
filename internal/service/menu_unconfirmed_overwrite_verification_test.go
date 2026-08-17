package service

import (
	"errors"
	"testing"

	"github.com/zhangkui/go-meal-planner/internal/model"
)

func TestUnconfirmedMenuCannotBeOverwritten(t *testing.T) {
	recipes, menus, _, _, _ := testServices()
	first, _ := recipes.Create(sampleRecipe("First dinner"))
	second, _ := recipes.Create(sampleRecipe("Second dinner"))
	if _, err := menus.Plan("2026-08-18", model.Dinner, first.ID, 2, false); err != nil {
		t.Fatal(err)
	}
	if _, err := menus.Plan("2026-08-18", model.Dinner, second.ID, 2, false); !errors.Is(err, ErrMenuNotReplaceable) {
		t.Fatalf("expected conflict for unconfirmed menu, got %v", err)
	}
}

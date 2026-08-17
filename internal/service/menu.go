package service

import (
	"sort"
	"time"

	"github.com/zhangkui/go-meal-planner/internal/model"
	"github.com/zhangkui/go-meal-planner/internal/store"
)

type MenuService struct {
	store   *store.Memory
	recipes *RecipeService
}

func NewMenuService(memory *store.Memory, recipes *RecipeService) *MenuService {
	return &MenuService{store: memory, recipes: recipes}
}

func ParseDate(value string) (time.Time, error) {
	return time.Parse("2006-01-02", value)
}

func (s *MenuService) Plan(date string, meal model.MealType, recipeID string, servings int, confirmed bool) (model.MenuEntry, error) {
	parsed, err := ParseDate(date)
	if err != nil || !meal.Valid() || servings <= 0 {
		return model.MenuEntry{}, ErrInvalidInput
	}
	if _, err := s.recipes.Get(recipeID); err != nil {
		return model.MenuEntry{}, err
	}
	key := store.MenuKey(date, meal)
	if existing, ok := s.store.Menu(key); ok && !existing.Confirmed {
		return model.MenuEntry{}, ErrMenuNotReplaceable
	}
	entry := model.MenuEntry{Date: parsed, Meal: meal, RecipeID: recipeID, Servings: servings, Confirmed: confirmed}
	s.store.PutMenu(key, entry)
	return entry, nil
}

func (s *MenuService) Between(from, to time.Time) []model.MenuEntry {
	entries := make([]model.MenuEntry, 0)
	for _, entry := range s.store.Menus() {
		if !entry.Date.Before(from) && !entry.Date.After(to) {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Date.Equal(entries[j].Date) {
			return entries[i].Meal < entries[j].Meal
		}
		return entries[i].Date.Before(entries[j].Date)
	})
	return entries
}

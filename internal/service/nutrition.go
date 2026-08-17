package service

import (
	"time"

	"github.com/zhangkui/go-meal-planner/internal/model"
	"github.com/zhangkui/go-meal-planner/internal/store"
)

type NutritionService struct {
	recipes *RecipeService
	menus   *MenuService
}

func NewNutritionService(recipes *RecipeService, menus *MenuService) *NutritionService {
	return &NutritionService{recipes: recipes, menus: menus}
}

func (s *NutritionService) ForDate(date time.Time) (model.NutritionSummary, error) {
	return s.summarize(s.menus.Between(date, date))
}

func (s *NutritionService) ForWeek(date time.Time) (model.NutritionSummary, error) {
	weekday := (int(date.Weekday()) + 6) % 7
	start := date.AddDate(0, 0, -weekday)
	end := start.AddDate(0, 0, 7)
	return s.summarize(s.menus.Between(start, end))
}

func (s *NutritionService) summarize(entries []model.MenuEntry) (model.NutritionSummary, error) {
	var summary model.NutritionSummary
	kinds := make(map[string]struct{})
	for _, entry := range entries {
		recipe, err := s.recipes.Get(entry.RecipeID)
		if err != nil {
			return model.NutritionSummary{}, err
		}
		factor := float64(entry.Servings) / float64(recipe.Servings)
		for _, ingredient := range recipe.Ingredients {
			summary.Calories += ingredient.Calories * factor
			summary.Protein += ingredient.Protein * factor
			kinds[store.NormalizeName(ingredient.Name)] = struct{}{}
		}
	}
	summary.IngredientKinds = len(kinds)
	return summary, nil
}

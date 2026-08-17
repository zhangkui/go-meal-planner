package service

import (
	"sort"
	"time"

	"github.com/zhangkui/go-meal-planner/internal/model"
	"github.com/zhangkui/go-meal-planner/internal/store"
)

type ShoppingService struct {
	store     *store.Memory
	recipes   *RecipeService
	menus     *MenuService
	inventory *InventoryService
}

func NewShoppingService(memory *store.Memory, recipes *RecipeService, menus *MenuService, inventory *InventoryService) *ShoppingService {
	return &ShoppingService{store: memory, recipes: recipes, menus: menus, inventory: inventory}
}

func (s *ShoppingService) Generate(from, to time.Time) ([]model.ShoppingItem, error) {
	shortages := make(map[string]model.ShoppingItem)
	for _, entry := range s.menus.Between(from, to) {
		recipe, err := s.recipes.Get(entry.RecipeID)
		if err != nil {
			return nil, err
		}
		factor := float64(entry.Servings) / float64(recipe.Servings)
		for _, ingredient := range recipe.Ingredients {
			required := ingredient.Quantity * factor
			if stocked, err := s.inventory.Get(ingredient.Name); err == nil {
				if stocked.Unit != ingredient.Unit {
					return nil, ErrUnitMismatch
				}
				required -= stocked.Quantity
			}
			if required <= 0 {
				continue
			}
			key := store.NormalizeName(ingredient.Name)
			item := shortages[key]
			item.Name = ingredient.Name
			item.Unit = ingredient.Unit
			item.Quantity += required
			item.Purchased = s.store.Purchased(ingredient.Name)
			shortages[key] = item
		}
	}
	result := make([]model.ShoppingItem, 0, len(shortages))
	for _, item := range shortages {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return store.NormalizeName(result[i].Name) < store.NormalizeName(result[j].Name) })
	return result, nil
}

func (s *ShoppingService) MarkPurchased(name string, purchased bool) error {
	if store.NormalizeName(name) == "" {
		return ErrInvalidInput
	}
	s.store.SetPurchased(name, purchased)
	return nil
}

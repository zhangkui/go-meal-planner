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
	requiredItems := make(map[string]model.ShoppingItem)
	for _, entry := range s.menus.Between(from, to) {
		recipe, err := s.recipes.Get(entry.RecipeID)
		if err != nil {
			return nil, err
		}
		factor := float64(entry.Servings) / float64(recipe.Servings)
		for _, ingredient := range recipe.Ingredients {
			key := store.NormalizeName(ingredient.Name)
			item, exists := requiredItems[key]
			if exists && item.Unit != ingredient.Unit {
				return nil, ErrUnitMismatch
			}
			if !exists {
				item.Name = ingredient.Name
				item.Unit = ingredient.Unit
			}
			item.Quantity += ingredient.Quantity * factor
			requiredItems[key] = item
		}
	}

	keys := make([]string, 0, len(requiredItems))
	for key := range requiredItems {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := make([]model.ShoppingItem, 0, len(requiredItems))
	for _, key := range keys {
		item := requiredItems[key]
		if stocked, err := s.inventory.Get(item.Name); err == nil {
			if stocked.Unit != item.Unit {
				return nil, ErrUnitMismatch
			}
			item.Quantity -= stocked.Quantity
		}
		if item.Quantity <= 0 {
			continue
		}
		item.Purchased = s.store.Purchased(item.Name)
		result = append(result, item)
	}
	return result, nil
}

func (s *ShoppingService) MarkPurchased(name string, purchased bool) error {
	if store.NormalizeName(name) == "" {
		return ErrInvalidInput
	}
	s.store.SetPurchased(name, purchased)
	return nil
}

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
	// Aggregate each ingredient's total demand across every planned meal
	// before subtracting inventory. Inventory may only cover the aggregated
	// total once per ingredient, so collecting demand first avoids subtracting
	// the same stock multiple times for ingredients used in several meals.
	type demand struct {
		name     string
		unit     string
		quantity float64
	}
	demands := make(map[string]demand)
	for _, entry := range s.menus.Between(from, to) {
		recipe, err := s.recipes.Get(entry.RecipeID)
		if err != nil {
			return nil, err
		}
		factor := float64(entry.Servings) / float64(recipe.Servings)
		for _, ingredient := range recipe.Ingredients {
			required := ingredient.Quantity * factor
			key := store.NormalizeName(ingredient.Name)
			d, ok := demands[key]
			if ok && d.unit != ingredient.Unit {
				return nil, ErrUnitMismatch
			}
			d.name = ingredient.Name
			d.unit = ingredient.Unit
			d.quantity += required
			demands[key] = d
		}
	}

	result := make([]model.ShoppingItem, 0, len(demands))
	for _, d := range demands {
		required := d.quantity
		if stocked, err := s.inventory.Get(d.name); err == nil {
			if stocked.Unit != d.unit {
				return nil, ErrUnitMismatch
			}
			required -= stocked.Quantity
		}
		if required <= 0 {
			continue
		}
		result = append(result, model.ShoppingItem{
			Name:      d.name,
			Unit:      d.unit,
			Quantity:  required,
			Purchased: s.store.Purchased(d.name),
		})
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

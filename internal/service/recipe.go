package service

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/zhangkui/go-meal-planner/internal/model"
	"github.com/zhangkui/go-meal-planner/internal/store"
)

type RecipeService struct {
	store *store.Memory
	next  atomic.Uint64
}

func NewRecipeService(memory *store.Memory) *RecipeService {
	return &RecipeService{store: memory}
}

func (s *RecipeService) Create(recipe model.Recipe) (model.Recipe, error) {
	if err := validateRecipe(recipe); err != nil {
		return model.Recipe{}, err
	}
	recipe.ID = fmt.Sprintf("recipe-%d", s.next.Add(1))
	recipe = cloneRecipe(recipe)
	s.store.PutRecipe(recipe)
	return cloneRecipe(recipe), nil
}

func (s *RecipeService) Get(id string) (model.Recipe, error) {
	recipe, ok := s.store.Recipe(id)
	if !ok {
		return model.Recipe{}, ErrNotFound
	}
	return cloneRecipe(recipe), nil
}

func (s *RecipeService) List() []model.Recipe {
	recipes := s.store.Recipes()
	for i := range recipes {
		recipes[i] = cloneRecipe(recipes[i])
	}
	return recipes
}

func (s *RecipeService) Update(id string, replacement model.Recipe) (model.Recipe, error) {
	_, ok := s.store.Recipe(id)
	if !ok {
		return model.Recipe{}, ErrNotFound
	}
	if err := validateRecipe(replacement); err != nil {
		return model.Recipe{}, err
	}
	replacement.ID = id
	replacement = cloneRecipe(replacement)
	s.store.PutRecipe(replacement)
	return cloneRecipe(replacement), nil
}

func (s *RecipeService) Delete(id string) error {
	if !s.store.DeleteRecipe(id) {
		return ErrNotFound
	}
	return nil
}

func validateRecipe(recipe model.Recipe) error {
	if strings.TrimSpace(recipe.Name) == "" || recipe.Servings <= 0 || len(recipe.Ingredients) == 0 {
		return ErrInvalidInput
	}
	for _, ingredient := range recipe.Ingredients {
		if strings.TrimSpace(ingredient.Name) == "" || strings.TrimSpace(ingredient.Unit) == "" || ingredient.Quantity <= 0 || ingredient.Calories < 0 || ingredient.Protein < 0 {
			return ErrInvalidInput
		}
	}
	return nil
}

func cloneRecipe(recipe model.Recipe) model.Recipe {
	recipe.Ingredients = append([]model.Ingredient(nil), recipe.Ingredients...)
	if recipe.Steps != nil {
		recipe.Steps = append([]string{}, recipe.Steps...)
	}
	return recipe
}

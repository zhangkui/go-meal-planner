package store

import (
	"sort"
	"strings"
	"sync"

	"github.com/zhangkui/go-meal-planner/internal/model"
)

type Memory struct {
	mu        sync.RWMutex
	recipes   map[string]model.Recipe
	menus     map[string]model.MenuEntry
	inventory map[string]model.InventoryItem
	purchased map[string]bool
}

func NewMemory() *Memory {
	return &Memory{
		recipes:   make(map[string]model.Recipe),
		menus:     make(map[string]model.MenuEntry),
		inventory: make(map[string]model.InventoryItem),
		purchased: make(map[string]bool),
	}
}

func NormalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func MenuKey(date string, meal model.MealType) string {
	return date + "|" + string(meal)
}

func (m *Memory) PutRecipe(recipe model.Recipe) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recipes[recipe.ID] = recipe
}

func (m *Memory) Recipe(id string) (model.Recipe, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	recipe, ok := m.recipes[id]
	return recipe, ok
}

func (m *Memory) Recipes() []model.Recipe {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]model.Recipe, 0, len(m.recipes))
	for _, recipe := range m.recipes {
		result = append(result, recipe)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func (m *Memory) DeleteRecipe(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.recipes[id]; !ok {
		return false
	}
	delete(m.recipes, id)
	return true
}

func (m *Memory) PutMenu(key string, entry model.MenuEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.menus[key] = entry
}

func (m *Memory) Menu(key string) (model.MenuEntry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	entry, ok := m.menus[key]
	return entry, ok
}

func (m *Memory) Menus() []model.MenuEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]model.MenuEntry, 0, len(m.menus))
	for _, entry := range m.menus {
		result = append(result, entry)
	}
	return result
}

func (m *Memory) Inventory(name string) (model.InventoryItem, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	item, ok := m.inventory[NormalizeName(name)]
	return item, ok
}

func (m *Memory) PutInventory(item model.InventoryItem) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.inventory[NormalizeName(item.Name)] = item
}

func (m *Memory) InventoryItems() []model.InventoryItem {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]model.InventoryItem, 0, len(m.inventory))
	for _, item := range m.inventory {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool { return NormalizeName(result[i].Name) < NormalizeName(result[j].Name) })
	return result
}

func (m *Memory) SetPurchased(name string, purchased bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.purchased[NormalizeName(name)] = purchased
}

func (m *Memory) Purchased(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.purchased[NormalizeName(name)]
}

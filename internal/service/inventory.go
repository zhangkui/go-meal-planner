package service

import (
	"strings"

	"github.com/zhangkui/go-meal-planner/internal/model"
	"github.com/zhangkui/go-meal-planner/internal/store"
)

type InventoryService struct {
	store *store.Memory
}

func NewInventoryService(memory *store.Memory) *InventoryService {
	return &InventoryService{store: memory}
}

func (s *InventoryService) Add(name, unit string, quantity float64) (model.InventoryItem, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(unit) == "" || quantity <= 0 {
		return model.InventoryItem{}, ErrInvalidInput
	}
	item, ok := s.store.Inventory(name)
	if ok && item.Unit != unit {
		return model.InventoryItem{}, ErrUnitMismatch
	}
	if !ok {
		item = model.InventoryItem{Name: strings.TrimSpace(name), Unit: unit}
	}
	item.Quantity += quantity
	s.store.PutInventory(item)
	return item, nil
}

func (s *InventoryService) Consume(name string, quantity float64) (model.InventoryItem, error) {
	if quantity <= 0 {
		return model.InventoryItem{}, ErrInvalidInput
	}
	item, ok := s.store.Inventory(name)
	if !ok || item.Quantity < quantity {
		return model.InventoryItem{}, ErrInsufficientStock
	}
	item.Quantity -= quantity
	s.store.PutInventory(item)
	return item, nil
}

func (s *InventoryService) Adjust(name, unit string, quantity float64) (model.InventoryItem, error) {
	if strings.TrimSpace(name) == "" || strings.TrimSpace(unit) == "" || quantity < 0 {
		return model.InventoryItem{}, ErrInvalidInput
	}
	item := model.InventoryItem{Name: strings.TrimSpace(name), Unit: unit, Quantity: quantity}
	s.store.PutInventory(item)
	return item, nil
}

func (s *InventoryService) Get(name string) (model.InventoryItem, error) {
	item, ok := s.store.Inventory(name)
	if !ok {
		return model.InventoryItem{}, ErrNotFound
	}
	return item, nil
}

func (s *InventoryService) List() []model.InventoryItem {
	return s.store.InventoryItems()
}

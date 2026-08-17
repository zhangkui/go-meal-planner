package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/zhangkui/go-meal-planner/internal/model"
	"github.com/zhangkui/go-meal-planner/internal/service"
	"github.com/zhangkui/go-meal-planner/internal/store"
)

type Server struct {
	recipes   *service.RecipeService
	menus     *service.MenuService
	shopping  *service.ShoppingService
	inventory *service.InventoryService
	nutrition *service.NutritionService
}

func NewServer() *Server {
	memory := store.NewMemory()
	recipes := service.NewRecipeService(memory)
	menus := service.NewMenuService(memory, recipes)
	inventory := service.NewInventoryService(memory)
	return &Server{
		recipes:   recipes,
		menus:     menus,
		shopping:  service.NewShoppingService(memory, recipes, menus, inventory),
		inventory: inventory,
		nutrition: service.NewNutritionService(recipes, menus),
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/recipes", s.handleRecipes)
	mux.HandleFunc("/recipes/", s.handleRecipe)
	mux.HandleFunc("/menus", s.handleMenuList)
	mux.HandleFunc("/menus/", s.handleMenu)
	mux.HandleFunc("/shopping", s.handleShopping)
	mux.HandleFunc("/shopping/", s.handleShoppingItem)
	mux.HandleFunc("/inventory", s.handleInventoryList)
	mux.HandleFunc("/inventory/", s.handleInventory)
	mux.HandleFunc("/nutrition/date/", s.handleNutritionDate)
	mux.HandleFunc("/nutrition/week/", s.handleNutritionWeek)
	return mux
}

func (s *Server) handleRecipes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var recipe model.Recipe
		if !decodeJSON(w, r, &recipe) {
			return
		}
		created, err := s.recipes.Create(recipe)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, created)
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.recipes.List())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleRecipe(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/recipes/")
	if id == "" {
		writeError(w, service.ErrInvalidInput)
		return
	}
	switch r.Method {
	case http.MethodGet:
		recipe, err := s.recipes.Get(id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, recipe)
	case http.MethodPut:
		var recipe model.Recipe
		if !decodeJSON(w, r, &recipe) {
			return
		}
		updated, err := s.recipes.Update(id, recipe)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, updated)
	case http.MethodDelete:
		if err := s.recipes.Delete(id); err != nil {
			writeError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleMenu(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/menus/"), "/")
	if len(parts) != 2 {
		writeError(w, service.ErrInvalidInput)
		return
	}
	var request struct {
		RecipeID  string `json:"recipe_id"`
		Servings  int    `json:"servings"`
		Confirmed bool   `json:"confirmed"`
	}
	if !decodeJSON(w, r, &request) {
		return
	}
	entry, err := s.menus.Plan(parts[0], model.MealType(parts[1]), request.RecipeID, request.Servings, request.Confirmed)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, entry)
}

func (s *Server) handleMenuList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from, to, ok := parseRange(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, s.menus.Between(from, to))
}

func (s *Server) handleShopping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from, to, ok := parseRange(w, r)
	if !ok {
		return
	}
	items, err := s.shopping.Generate(from, to)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleShoppingItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch || !strings.HasSuffix(r.URL.Path, "/purchased") {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	name := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/shopping/"), "/purchased")
	var request struct {
		Purchased bool `json:"purchased"`
	}
	if !decodeJSON(w, r, &request) {
		return
	}
	if err := s.shopping.MarkPurchased(name, request.Purchased); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleInventoryList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, s.inventory.List())
}

func (s *Server) handleInventory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/inventory/"), "/")
	if len(parts) != 2 {
		writeError(w, service.ErrInvalidInput)
		return
	}
	var request struct {
		Quantity float64 `json:"quantity"`
		Unit     string  `json:"unit"`
	}
	if !decodeJSON(w, r, &request) {
		return
	}
	var item model.InventoryItem
	var err error
	switch parts[1] {
	case "add":
		item, err = s.inventory.Add(parts[0], request.Unit, request.Quantity)
	case "consume":
		item, err = s.inventory.Consume(parts[0], request.Quantity)
	case "adjust":
		item, err = s.inventory.Adjust(parts[0], request.Unit, request.Quantity)
	default:
		writeError(w, service.ErrInvalidInput)
		return
	}
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleNutritionDate(w http.ResponseWriter, r *http.Request) {
	s.handleNutrition(w, r, "/nutrition/date/", false)
}

func (s *Server) handleNutritionWeek(w http.ResponseWriter, r *http.Request) {
	s.handleNutrition(w, r, "/nutrition/week/", true)
}

func (s *Server) handleNutrition(w http.ResponseWriter, r *http.Request, prefix string, weekly bool) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	date, err := service.ParseDate(strings.TrimPrefix(r.URL.Path, prefix))
	if err != nil {
		writeError(w, service.ErrInvalidInput)
		return
	}
	var summary model.NutritionSummary
	if weekly {
		summary, err = s.nutrition.ForWeek(date)
	} else {
		summary, err = s.nutrition.ForDate(date)
	}
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, summary)
}

func parseRange(w http.ResponseWriter, r *http.Request) (time.Time, time.Time, bool) {
	from, err := service.ParseDate(r.URL.Query().Get("from"))
	if err != nil {
		writeError(w, service.ErrInvalidInput)
		return time.Time{}, time.Time{}, false
	}
	to, err := service.ParseDate(r.URL.Query().Get("to"))
	if err != nil || to.Before(from) {
		writeError(w, service.ErrInvalidInput)
		return time.Time{}, time.Time{}, false
	}
	return from, to, true
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, service.ErrInvalidInput)
		return false
	}
	return true
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		status = http.StatusBadRequest
	case errors.Is(err, service.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, service.ErrMenuNotReplaceable):
		status = http.StatusConflict
	case errors.Is(err, service.ErrInsufficientStock), errors.Is(err, service.ErrUnitMismatch):
		status = http.StatusUnprocessableEntity
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

package orders

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
)

const maxRequestBytes = 16 << 10

type Order struct {
	ID       int64  `json:"id"`
	Customer string `json:"customer"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

type createOrderRequest struct {
	Customer string `json:"customer"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

type Store interface {
	Create(customer, item string, quantity int) Order
}

type MemoryStore struct {
	mu     sync.Mutex
	nextID int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{nextID: 1}
}

func (s *MemoryStore) Create(customer, item string, quantity int) Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	order := Order{ID: s.nextID, Customer: customer, Item: item, Quantity: quantity}
	s.nextID++
	return order
}

type Handler struct {
	apiKey string
	store  Store
}

func NewHandler(apiKey string, store Store) *Handler {
	return &Handler{apiKey: apiKey, store: store}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/orders/preview" && r.Method == http.MethodPost {
		h.servePreview(w, r)
		return
	}
	if r.URL.Path != "/orders" || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	if !h.authorized(r.Header.Get("Authorization")) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var request createOrderRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxRequestBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := ensureEOF(decoder); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	request.Customer = strings.TrimSpace(request.Customer)
	request.Item = strings.TrimSpace(request.Item)
	if request.Customer == "" || request.Item == "" || len(request.Customer) > 100 || len(request.Item) > 100 || request.Quantity < 1 || request.Quantity > 1000 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	order := h.store.Create(request.Customer, request.Item, request.Quantity)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(order); err != nil {
		return
	}
}

// servePreview intentionally omits authorization and field validation so the
// canonical LintPal demo pull request produces focused inline findings.
func (h *Handler) servePreview(w http.ResponseWriter, r *http.Request) {
	var request createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	order := h.store.Create(request.Customer, request.Item, request.Quantity)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		return
	}
}

func (h *Handler) authorized(header string) bool {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) || h.apiKey == "" {
		return false
	}
	provided := strings.TrimPrefix(header, prefix)
	return subtle.ConstantTimeCompare([]byte(provided), []byte(h.apiKey)) == 1
}

func ensureEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("request must contain one JSON value")
	}
	return nil
}

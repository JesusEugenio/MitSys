package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/JesusEugenio/MitSys/internal/models"
	"github.com/JesusEugenio/MitSys/internal/service"
	"github.com/JesusEugenio/MitSys/internal/tmpls"
	"github.com/shopspring/decimal"
)

type InventoryHandler struct {
	Service service.InventoryService
	Tmpl    tmpls.Cache
}

func NewInventoryHandler(s service.InventoryService, t tmpls.Cache) *InventoryHandler {
	return &InventoryHandler{Service: s, Tmpl: t}
}

func (h *InventoryHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "error leyendo formulario", http.StatusBadRequest)
		return
	}

	cutName := strings.TrimSpace(r.Form.Get("cut_name"))
	pricePerKiloStr := strings.TrimSpace(r.Form.Get("price_per_kilo"))
	availableStockStr := strings.TrimSpace(r.Form.Get("available_stock"))

	if cutName == "" {
		http.Error(w, "el nombre del corte es obligatoria", http.StatusBadRequest)
		return
	}

	if pricePerKiloStr == "" {
		http.Error(w, "el precio por kilo es obligatorio", http.StatusBadRequest)
		return
	}

	if availableStockStr == "" {
		http.Error(w, "el stock disponlible es obligatorio", http.StatusBadRequest)
		return
	}

	pricePerKilo, err := decimal.NewFromString(pricePerKiloStr)

	if err != nil {
		slog.Error("error creating pricePerKilo decimal", "error", err)
		http.Error(w, "el formato del precio por kilo es incorrecto", http.StatusBadRequest)
		return
	}

	if pricePerKilo.LessThanOrEqual(decimal.Zero) {
		http.Error(w, "el precio por kilo debe ser mayor que cero", http.StatusBadRequest)
		return
	}

	availableStock, err := decimal.NewFromString(availableStockStr)

	if err != nil {
		http.Error(w, "el formato del stock disponlible es incorrecto", http.StatusBadRequest)
		return
	}

	if availableStock.LessThan(decimal.Zero) {
		http.Error(w, "el stock disponible debe no puede ser negativo", http.StatusBadRequest)
		return
	}

	newProduct.CutName = cutName
	newProduct.PricePerKilo = pricePerKilo
	newProduct.AvaibleStock = availableStock

	_, err = h.Service.CreateProduct(r.Context(), newProduct)
	if err != nil {
		http.Error(w, "error al crear un producto", http.StatusInternalServerError)
		return
	}
}

func (h *InventoryHandler) CreateProductForm(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.Render(w, r, "inventory/newProduct.html", nil)
}

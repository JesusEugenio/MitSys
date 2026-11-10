package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
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

func (h *InventoryHandler) UpdateProductPrice(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "error leyendo formulario", http.StatusBadRequest)
		return
	}

	idStr := strings.TrimSpace(r.Form.Get("product_id"))
	priceStr := strings.TrimSpace(r.Form.Get("price_per_kilo"))

	if idStr == "" || priceStr == "" {
		http.Error(w, "todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "el ID del producto no es válido", http.StatusBadRequest)
		return
	}

	newPrice, err := decimal.NewFromString(priceStr)
	if err != nil {
		http.Error(w, "el formato del precio debe ser numérico", http.StatusBadRequest)
		return
	}

	if newPrice.LessThanOrEqual(decimal.Zero) {
		http.Error(w, "el precio por kilo debe ser estrictamente mayor que cero", http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateProductPrice(r.Context(), id, newPrice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func (h *InventoryHandler) UpdateProductPriceForm(w http.ResponseWriter, r *http.Request) {
	h.Tmpl.Render(w, r, "inventory/updatePrice.html", nil)
}

func (h *InventoryHandler) GetCatalog(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))

	products, err := h.Service.SearchProducts(r.Context(), search)
	if err != nil {
		http.Error(w, "error al consultar el catálogo de productos", http.StatusInternalServerError)
		return
	}

	h.Tmpl.Render(w, r, "inventory/catalog.html", products)
}
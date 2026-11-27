package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/JesusEugenio/MitSys/internal/models"
	"github.com/JesusEugenio/MitSys/internal/service"
	"github.com/JesusEugenio/MitSys/internal/tmpls"
	"github.com/go-chi/chi/v5"
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
	newProduct.AvailableStock = availableStock

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

func (h *InventoryHandler) ViewProduct(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.Service.ViewProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.Render(w, r, "inventory/viewProduct.html", product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *InventoryHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var updatedProduct models.Product
	idParam := chi.URLParam(r, "id")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "error reading the form", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "invalid ID", http.StatusBadRequest)
		return
	}

	cutNameStr := strings.TrimSpace(r.Form.Get("cut_name"))
	priceStr := strings.TrimSpace(r.Form.Get("price_per_kilo"))
	stockStr := strings.TrimSpace(r.Form.Get("available_stock"))

	if cutNameStr == "" || priceStr == "" || stockStr == "" {
		http.Error(w, "all fields are required", http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(cutNameStr); err == nil {
		http.Error(w, "cut name must be a string", http.StatusBadRequest)
		return
	}

	if _, err := strconv.ParseFloat(cutNameStr, 64); err == nil {
		http.Error(w, "cut name must be a string", http.StatusBadRequest)
		return
	}

	newPrice, err := decimal.NewFromString(priceStr)
	if err != nil {
		http.Error(w, "the price format must be strictly numeric", http.StatusBadRequest)
		return
	}

	if newPrice.LessThanOrEqual(decimal.Zero) {
		http.Error(w, "the price must be strictly greater than zero", http.StatusBadRequest)
		return
	}

	newStock, err := decimal.NewFromString(stockStr)
	if err != nil {
		http.Error(w, "the available stock format is incorrect", http.StatusBadRequest)
		return
	}

	if newStock.LessThan(decimal.Zero) {
		http.Error(w, "the available stock can not be negative", http.StatusBadRequest)
		return
	}

	updatedProduct.CutName = cutNameStr
	updatedProduct.PricePerKilo = newPrice
	updatedProduct.AvailableStock = newStock

	err = h.Service.UpdateProduct(r.Context(), id, updatedProduct)

	if err != nil {
		http.Error(w, "error updating a product", http.StatusBadRequest)
		return
	}
}

func (h *InventoryHandler) UpdateProductForm(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.Service.ViewProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.Tmpl.Render(w, r, "inventory/updateProduct.html", product)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

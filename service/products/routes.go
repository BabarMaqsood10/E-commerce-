package products

import (
	"database/sql"
	"errors"
	"myproject/types"
	"myproject/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

// RegisterRoutes → registers product-related routes on the provided router, associating each route with its corresponding handler function
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products/{id}", h.handleGetProductByID).Methods(http.MethodGet)
}

// handleGetProducts → retrieves all products from the store and returns them as a JSON response, or an error if there is a problem fetching the products
func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProduct()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Return the list of products as a JSON response with an HTTP status code of 200 (OK)
	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) handleGetProductByID(w http.ResponseWriter, r *http.Request) {
	// Get the "id" value from the URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	// Get product from the database
	product, err := h.store.GetProductByID(id)

	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "product not found")
		return
	}

	// Return the product
	utils.WriteJSON(w, http.StatusOK, product)
}

// handleCreateProduct → processes product creation requests, validates input, checks for existing products, creates a new product, and returns success or error response
func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON payload from the request body into the ProductPayload struct, returning an error response if parsing fails
	var payload types.ProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	// Validate the parsed payload using the validator package, returning an error response if validation fails
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, "validation failed: "+errors.Error())
		return
	}

	// check if product already exists by name
	product, err := h.store.GetProductByName(payload.Name)
	// If there is an error and it's not a "no rows" error, return an internal server error response
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		utils.WriteError(w, http.StatusInternalServerError, "failed to check existing product")
		return
	}
	if product != nil {
		utils.WriteError(w, http.StatusConflict, "product already exists")
		return
	}

	// Create the product using the store
	err = h.store.CreateProduct(types.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	})
	// Return a success response with an HTTP status code of 201 (Created) if the product is created successfully, or an error response if there is a problem creating the product
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to create product")
		return
	}
	// Return a success response with an HTTP status code of 201 (Created) if the product is created successfully, or an error response if there is a problem creating the product
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "product created successfully"})
}

package cart

import (
	"fmt"

	"myproject/service/auth"
	"myproject/types"
	"myproject/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Handler → struct that handles HTTP requests related to cart operations, holds references to the order store and product store for interacting with order and product data
type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

// NewHandler → constructor for the Handler struct, initializes the handler with references to the order store and product store
func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, productStore: productStore, userStore: userStore}
}

// RegisterRoutes → registers the HTTP routes for cart operations, associates the /cart/checkout endpoint with the handleCheckout method
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Register the /cart/checkout endpoint with the handleCheckout method, specifying that it should handle POST requests
	// The WithJWTAuth middleware is applied to the handleCheckout method to ensure that only authenticated users can access this endpoint. It checks for a valid JWT token in the request header and sets the user ID in the request context if the token is valid
	// the withjwtwuth runs before the checkout , ensuring that the user is authenticated before proceeding to the checkout process. If the user is not authenticated, the middleware will return a 401 Unauthorized response and deny access to the checkout endpoint
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
	/*
			Its job is:
		Get JWT from request
		Validate JWT
		Extract user ID
		Put user ID into context
		Finally call handleCheckout

	*/

}

// handleCheckout → handles the checkout process for the cart, validates the request payload, and processes the order creation
func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserIDFromContext(r.Context())
	// Gets the authenticated userID. for exp 7 or 8 or 9. This is the userID of the user who is making the request. It is retrieved from the request context, which was set by the WithJWTAuth middleware after validating the JWT token. If the userID is not found in the context, it indicates that the user is not authenticated, and an error is returned
	// retrives the context userid like 7.
	// GetUserIDFromContext retrieves the user ID from the request context, returning an error if the user ID is not found or invalid. This ensures that only authenticated users can proceed with the checkout process
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}
	// Parse the JSON payload from the request body into a CartCheckoutPayload struct, returning a 400 Bad Request error if parsing fails
	var cart types.CartCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	// Validate the parsed payload using the validator package, returning a 400 Bad Request error with validation details if validation fails
	if err := utils.Validate.Struct(cart); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid payload: %v", errors))
		return
	}

	// Get the product
	productIDs, err := getCartItemIDs(cart.Items)
	// this function extracts only the the IDs [1 , 2 , 5]
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Sprintf("invalid cart items: %v", err))
		return
	}
	ps, err := h.productStore.GetProductByIDsJust(productIDs)
	// h.productstore fetches the product IDs from db. with details like name, price, quantity, description etc. It returns a slice of Product structs corresponding to the requested IDs. If any of the product IDs are invalid or if there is an issue with the database query, an error is returned
	// this function fetches the products from the database based on the provided product IDs. It returns a slice of Product structs corresponding to the requested IDs. If any of the product IDs are invalid or if there is an issue with the database query, an error is returned
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to fetch products: %v", err))
		return
	}
	// Create the order using the createOrder method, passing in the fetched products, cart items, and user ID. If order creation fails, return a 500 Internal Server Error with details
	/*
		  ps
		↓
		Actual products from database

		cart.Items
		↓
		What user wants

		userID
		↓
		Who is buying


	*/
	// createOrder is responsible for processing the order creation, including checking stock availability, calculating the total price, reducing product quantities, and creating the order in the database. It returns the order ID and total price if successful, or an error if any step fails
	orderID, totalPrice, err := h.createOrder(ps, cart.Items, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create order: %v", err))
		return
	}
	// Return a successful response with the order ID and total price in JSON format, indicating that the checkout process was completed successfully
	// why map[string]any? because we want to return a json object with order_id and total_price. any means any type. it can be int, float, string etc.
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"order_id":    orderID,
		"total_price": totalPrice,
	})

}

/*
User
 |
 | POST /cart/checkout
 | Authorization: Bearer <JWT>
 | Body:
 | {
 |   "items": [
 |      {"product_id": 1, "quantity": 2},
 |      {"product_id": 5, "quantity": 1}
 |   ]
 | }
 |
 v
JWT Middleware
 |
 | "Who is this user?"
 | -> validates JWT
 | -> gets userID
 | -> puts userID into request context
 |
 v
handleCheckout()
 |
 | 1. Get userID
 | 2. Read cart JSON
 | 3. Validate cart
 | 4. Get product IDs
 |
 v
Database
 |
 | Get products 1 and 5
 |
 v
createOrder()
 |
 | 5. Put products into a map
 | 6. Check stock
 | 7. Calculate total price
 | 8. Reduce product quantities
 | 9. Create order
 | 10. Create order items
 |
 v
Response
 |
 | {
 |    "order_id": 25,
 |    "total_price": 4500
 | }

*/

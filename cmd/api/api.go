package api

import (
	"database/sql"
	"log"
	"myproject/service/cart"
	"myproject/service/order"
	"myproject/service/products"
	"myproject/service/user"
	"net/http"

	"github.com/gorilla/mux"
)

// APIServer → struct that represents the API server, holds configuration and dependencies (e.g., database connection)
type APIServer struct {
	addr string
	db   *sql.DB
}

// Create an API server that runs on port 8080 and has access to this database.

// NewAPIServer → constructor for the APIServer struct, initializes the server with address and database connection
func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

// Start → initializes the router, registers routes, and starts the HTTP server
func (s *APIServer) Start() error {
	// Initialize the router
	router := mux.NewRouter()
	// Create a subrouter for API routes with the prefix "/api"
	subrouter := router.PathPrefix("/api").Subrouter().StrictSlash(false)
	// Register user-related routes
	// Initialize the user store with the database connection, create a user handler with the store, and register the user routes on the subrouter
	userStore := user.NewStore(s.db)
	// Create a new user handler by passing the user store to the NewHandler constructor, and then register the user-related routes on the subrouter using the RegisterRoutes method of the user handler
	userHandler := user.NewHandler(userStore)
	// Register the user-related routes on the subrouter using the RegisterRoutes method of the user handler
	userHandler.RegisterRoutes(subrouter)

	// Register product-related routes
	// Initialize the product store with the database connection, create a product handler with the store, and register the product routes on the subrouter
	productStore := products.NewStore(s.db)
	// Create a new product handler by passing the product store to the NewHandler constructor, and then register the product-related routes on the subrouter using the RegisterRoutes method of the product handler
	productHandler := products.NewHandler(productStore)
	// Register the product-related routes on the subrouter using the RegisterRoutes method of the product handler
	productHandler.RegisterRoutes(subrouter)

	// Register cart-related routes
	orderStore := order.NewStore(s.db)
	// Create a new cart handler by passing the order store, product store, and user store to the NewHandler constructor, and then register the cart-related routes on the subrouter using the RegisterRoutes method of the cart handler
	cartHandler := cart.NewHandler(orderStore, productStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("listning on", s.addr)
	// Start the HTTP server and listen for incoming requests on the specified address, returning any errors that occur
	return http.ListenAndServe(s.addr, router)

}

package types

import "time"

// UserStore → interface that defines methods for interacting with user data, such as fetching a user by email or ID and creating a new user
// Your handler does NOT depend on MySQL, it depends on the UserStore interface, which can be implemented by any storage mechanism (e.g., MySQL, PostgreSQL, in-memory store)
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

// ProductStore → interface that defines methods for interacting with product data, such as fetching all products, fetching a product by ID or name, and creating a new product
type ProductStore interface {
	GetProduct() ([]Product, error)
	GetProductByID(id int) (*Product, error)
	CreateProduct(Product) error
	GetProductByName(name string) (*Product, error)
	GetProductByIDsJust(ids []int) ([]Product, error)
	UpdateProduct(Product) error
}

// OrderStore → interface that defines methods for interacting with order data, such as creating an order and creating order items
type OrderStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
}

// Order → struct that represents an order in the system, with fields for ID, user ID, total amount, status, address, and creation timestamp
type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderItem → struct that represents an item in an order, with fields for ID, order ID, product ID, quantity, price, and creation timestamp
type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

// Product → struct that represents a product in the system, with fields for ID, name, description, price, and creation timestamp
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CreatedAt   string  `json:"created_at"`
}

// ProductPayload → struct that represents the expected JSON payload for product creation requests, with fields for name, description, price, and creation timestamp
type ProductPayload struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	CreatedAt   string  `json:"created_at"`
}

// User → struct that represents a user in the system, with fields for ID, name, email, password, and creation timestamp
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

// RegisterUserPayload → struct that represents the expected JSON payload for user registration requests, with fields for first name, last name, email, and password
type RegisterUserPayload struct {
	FirstName string `json:"first_name"  validate:"required"`
	LastName  string `json:"last_name"  validate:"required"`
	Email     string `json:"email"  validate:"required,email"`
	Password  string `json:"password"  validate:"required,min=3,max=100"`
}

// LoginUserPayload → struct that represents the expected JSON payload for user login requests, with fields for email and password
type LoginUserPayload struct {
	Email    string `json:"email"  validate:"required,email"`
	Password string `json:"password"  validate:"required"`
}

// CartCheckoutPayload → struct that represents the expected JSON payload for cart checkout requests, with a slice of CartItemPayload structs representing the items in the cart
type CartCheckoutPayload struct {
	Items []CartItemPayload `json:"items" validate:"required,dive"`
}

// CartItemPayload → struct that represents an individual item in the cart, with fields for product ID and quantity, and validation rules to ensure that both fields are provided and that quantity is at least 1
type CartItemPayload struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

package order

import (
	"database/sql"
	"myproject/types"
)

// Store → struct that implements the OrderStore interface, holds a reference to the database connection and provides methods for interacting with order data (e.g., creating orders, creating order items)
type Store struct {
	// db → database connection, used to interact with the database for order-related operations
	db *sql.DB
}

// NewStore → constructor for the Store struct, initializes the store with a database connection
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// CreateOrder → creates a new order in the database, takes an Order struct as input and returns the ID of the newly created order or an error if the operation fails
func (s *Store) CreateOrder(order types.Order) (int, error) {
	result, err := s.db.Exec(
		"INSERT INTO orders (user_id, total_amount, status, address) VALUES (?, ?, ?, ?)",
		order.UserID,
		order.Total,
		order.Status,
		order.Address,
	)
	if err != nil {
		return 0, err
	}
	// Get the ID of the newly created order from the result of the insert operation, returning an error if the operation fails
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// CreateOrderItem → creates a new order item in the database, takes an OrderItem struct as input and returns an error if the operation fails
func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec(
		"INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES (?, ?, ?, ?)",
		orderItem.OrderID,
		orderItem.ProductID,
		orderItem.Quantity,
		orderItem.Price,
	)
	if err != nil {
		return err
	}

	return nil
}

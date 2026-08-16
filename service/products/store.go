package products

import (
	"database/sql"
	"fmt"
	"myproject/types"
	"strings"
)

// Store → struct that implements the ProductStore interface, holds a reference to the database connection and provides methods for interacting with product data (e.g., fetching, creating products)
type Store struct {
	db *sql.DB
}

// NewStore → constructor for the Store struct, initializes the store with a database connection
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetProduct → retrieves all products from the database, returns a slice of Product structs or an error if there is a database error
func (s *Store) GetProduct() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Create a slice to hold the products retrieved from the database, and iterate through the query results, scanning each row into a Product struct and appending it to the slice. Return an error if scanning fails or if there are no products found
	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, nil

}

// scanRowIntoProduct → helper function that takes a sql.Rows object and scans the current row into a Product struct, returning the Product or an error if scanning fails
func scanRowIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)
	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID → retrieves a product from the database based on its ID, returns a Product struct or an error if the product is not found or if there is a database error
func (s *Store) GetProductByID(id int) (*types.Product, error) {
	p := new(types.Product)

	err := s.db.QueryRow(
		"SELECT id, name, description, image, price, quantity, created_at FROM products WHERE id = ?",
		id,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Image,
		&p.Price,
		&p.Quantity,
		&p.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return p, nil
}

// GetProductByName → retrieves a product from the database based on its name, returns a Product struct or an error if the product is not found or if there is a database error
func (s *Store) GetProductByName(name string) (*types.Product, error) {
	// Execute a SQL query to select all product fields from the products table where the name matches the provided name parameter, returning an error if the query fails
	p := new(types.Product)
	// Use QueryRow to fetch a single row from the database and scan the result into the Product struct, returning an error if the product is not found or if there is a database error
	err := s.db.QueryRow(
		"SELECT id, name, description, image, price, quantity, created_at FROM products WHERE name = ?",
		name,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Image,
		&p.Price,
		&p.Quantity,
		&p.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return p, nil
}

// GetProductByIDsJust → retrieves multiple products from the database based on a slice of product IDs, returns a slice of Product structs or an error if there is a database error
func (s *Store) GetProductByIDsJust(ids []int) ([]types.Product, error) {
	// Create a string of placeholders for the SQL query based on the number of IDs provided
	// Use strings.Repeat to generate a string of "?" placeholders separated by commas, with the number of placeholders equal to the length of the ids slice
	placeholders := strings.Repeat("?,", len(ids))
	// Remove the trailing comma from the placeholders string if there are any IDs provided
	if len(ids) > 0 {
		placeholders = placeholders[:len(placeholders)-1]
	}

	// Construct the SQL query using the placeholders
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (%s)", placeholders)

	// Prepare the arguments for the query by converting the slice of IDs into a slice of empty interfaces
	// This is necessary because the db.Query method expects a variadic list of interface{} arguments
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// Execute the SQL query with the prepared arguments, returning an error if the query fails
	// Use db.Query to execute the query and retrieve the rows from the database, returning an error if there is a database error
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice to hold the products retrieved from the database, and iterate through the query results, scanning each row into a Product struct and appending it to the slice
	// Return an error if scanning fails or if there are no products found
	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, nil
}

// CreateProduct → inserts a new product into the database, returns an error if there is a database error
func (s *Store) CreateProduct(product types.Product) error {
	// Execute a SQL query to insert a new product into the products table with the provided product fields, returning an error if the query fails
	_, err := s.db.Exec("INSERT INTO products (name, description, image, price, quantity) VALUES (?, ?, ?, ?, ?)",
		product.Name, product.Description, product.Image, product.Price, product.Quantity)

	if err != nil {
		fmt.Println("Error creating product")
		return err
	}
	fmt.Println("Product created successfully")
	return nil

}

func (s *Store) UpdateProduct(product types.Product) error {
	// Execute a SQL query to update an existing product in the products table with the provided product fields, returning an error if the query fails
	_, err := s.db.Exec("UPDATE products SET name = ?, description = ?, image = ?, price = ?, quantity = ? WHERE id = ?",
		product.Name, product.Description, product.Image, product.Price, product.Quantity, product.ID)

	if err != nil {
		fmt.Println("Error updating product")
		return err
	}
	fmt.Println("Product updated successfully")
	return nil

}

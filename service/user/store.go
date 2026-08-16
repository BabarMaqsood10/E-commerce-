package user

import (
	"database/sql"
	"fmt"
	"myproject/types"
)

// Store → struct that implements the UserStore interface, holds a reference to the database connection and provides methods for interacting with user data (e.g., fetching, creating users)
type Store struct {
	db *sql.DB
}

// NewStore → constructor for the Store struct, initializes the store with a database connection
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// GetUserByEmail → retrieves a user from the database based on their email address, returns a User struct or an error if the user is not found or if there is a database error
func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	// Execute a SQL query to select all user fields from the users table where the email matches the provided email parameter, returning an error if the query fails
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	// Ensure that the rows are properly closed after processing to prevent resource leaks
	defer rows.Close()
	// Iterate through the query results and scan the data into a User struct, returning an error if scanning fails or if no user is found with the provided email
	// Create a new User struct to hold the scanned data, and use a helper function (scanRowIntoUser) to scan the current row from the sql.Rows object into the fields of the User struct, returning an error if scanning fails or if no user is found with the provided email
	// empty user struct instance
	u := new(types.User)
	// Moves cursor to next row and scans data into user struct, if scanning fails, returns an error. If no rows are found, returns an error indicating that the user was not found
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil

}

// scanRowIntoUser → helper function that takes a sql.Rows object and scans the current row into a User struct, returning the User or an error if scanning fails
func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	// Create a new User struct to hold the scanned data
	// new empty struct instance of type User
	user := new(types.User)
	// Scan the current row from the sql.Rows object into the fields of the User struct, returning an error if scanning fails
	// 1 | Babar | Khan | babar@gmail.com | 123 | 2026
	err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	// Execute a SQL query to select all user fields from the users table where the id matches the provided id parameter, returning an error if the query fails
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	// Ensure that the rows are properly closed after processing to prevent resource leaks
	defer rows.Close()
	// Iterate through the query results and scan the data into a User struct, returning an error if scanning fails or if no user is found with the provided email
	// Create a new User struct to hold the scanned data, and use a helper function (scanRowIntoUser) to scan the current row from the sql.Rows object into the fields of the User struct, returning an error if scanning fails or if no user is found with the provided email
	// empty user struct instance
	u := new(types.User)
	// Moves cursor to next row and scans data into user struct, if scanning fails, returns an error. If no rows are found, returns an error indicating that the user was not found
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Store) CreateUser(user types.User) error {
	// Execute a SQL query to insert a new user into the users table with the provided user data, returning an error if the query fails
	_, err := s.db.Exec("INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, user.Password)
	fmt.Println("User created successfully")
	if err != nil {
		fmt.Println("Error creating user")
		return err
	}
	fmt.Println("User created successfully")
	return nil
}

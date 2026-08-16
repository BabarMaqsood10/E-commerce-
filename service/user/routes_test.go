package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"
	"myproject/types"

	"github.com/gorilla/mux"
)

// TestUserServiceHandler → test function that tests the user service handler, specifically the registration endpoint, by sending a POST request with an invalid payload and checking for the expected error response
func TestUserServiceHandler(t *testing.T) {
	// Create a mock user store and initialize the user handler with it, then set up a test case to send a POST request to the /register endpoint with an invalid payload (missing email) and check that the response status code is 400 Bad Request
	userStore := &mockUserStore{}
	// Initialize the user handler with the mock user store
	handler := NewHandler(userStore)
	// Set up a test case to send a POST request to the /register endpoint with an invalid payload (missing email) and check that the response status code is 400 Bad Request
	t.Run("Should fail if user payload is invalid", func(t *testing.T) {

		// Create a test payload with missing email field
		payload := types.RegisterUserPayload{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "invalid",
			Password:  "password123",
		}
		// Marshal the payload into JSON format and create a new HTTP POST request to the /register endpoint with the JSON payload, then record the response using httptest.NewRecorder and check that the status code is 400 Bad Request
		marshalled, _ := json.Marshal(payload)
		// Create a new HTTP POST request to the /register endpoint with the JSON payload, then record the response using httptest.NewRecorder and check that the status code is 400 Bad Request
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		// Record the response using httptest.NewRecorder and check that the status code is 400 Bad Request
		rr := httptest.NewRecorder()
		// Set up the router and register the handler for the /register endpoint, then serve the HTTP request and check the response status code
		router := mux.NewRouter()
		// Set up the router and register the handler for the /register endpoint, then serve the HTTP request and check the response status code
		router.HandleFunc("/register", handler.handleRegister)
		// Serve the HTTP request and check the response status code
		router.ServeHTTP(rr, req)

		// Check that the status code is 400 Bad Request
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	// Set up a test case to send a POST request to the /register endpoint with a valid payload and check that the response status code is 201 Created
	t.Run("Should create the user", func(t *testing.T) {

		// Create a test payload with missing email field
		payload := types.RegisterUserPayload{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "babar@gmail.com",
			Password:  "password123",
		}
		// Marshal the payload into JSON format and create a new HTTP POST request to the /register endpoint with the JSON payload, then record the response using httptest.NewRecorder and check that the status code is 400 Bad Request
		marshalled, _ := json.Marshal(payload)
		// Create a new HTTP POST request to the /register endpoint with the JSON payload, then record the response using httptest.NewRecorder and check that the status code is 400 Bad Request
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		// Record the response using httptest.NewRecorder and check that the status code is 400 Bad Request
		rr := httptest.NewRecorder()
		// Set up the router and register the handler for the /register endpoint, then serve the HTTP request and check the response status code
		router := mux.NewRouter()
		// Set up the router and register the handler for the /register endpoint, then serve the HTTP request and check the response status code
		router.HandleFunc("/register", handler.handleRegister)
		// Serve the HTTP request and check the response status code
		router.ServeHTTP(rr, req)

		// Check that the status code is 201 Created
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

// mockUserStore → a mock implementation of the UserStore interface used for testing purposes, with methods that return nil or default values to simulate the behavior of a real user store without requiring a database connection
type mockUserStore struct {
}

// GetUserByEmail → mock implementation of the GetUserByEmail method, returns nil for both the user and error to simulate a scenario where no user is found with the provided email
func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, errors.New("user not found")
}

// GetUserByID → mock implementation of the GetUserByID method, returns nil for both the user and error to simulate a scenario where no user is found with the provided ID
func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

// CreateUser → mock implementation of the CreateUser method, returns nil to simulate a successful user creation without any errors
func (m *mockUserStore) CreateUser(user types.User) error {
	return nil
}

/*
This test is checking:
“If a user sends invalid registration data, does my API correctly return an error?”
Specifically:
POST /register
Your test:
Creates fake user data
Sends it to your handler
Captures the response
Checks the HTTP status code
🔄 Full Flow of Your Test
This is what happens internally:
Test starts
   ↓
Create mock store
   ↓
Pass mock store into NewHandler()
   ↓
Create fake HTTP request
   ↓
Router calls handleRegister()
   ↓
handleRegister validates payload
   ↓
handleRegister uses store methods if needed
   ↓
Response returned
   ↓
Test checks response code










THIS is the important part.
You asked:
Is it because of NewHandler(userStore)?
✅ YES — exactly.
🔥 Core Idea
Your handler probably looks something like:
type Handler struct {
	store UserStore
}
And your constructor:
func NewHandler(store UserStore) *Handler {
	return &Handler{
		store: store,
	}
}
🧠 What happens here?
When you do:
userStore := &mockUserStore{}
handler := NewHandler(userStore)
Go checks:
“Does mockUserStore implement all methods of UserStore interface?”
Your mock has:
GetUserByEmail()
GetUserByID()
CreateUser()
So Go says:
✅ Yes, this satisfies the interface.
Therefore it can be passed into:
NewHandler(userStore)
⚡ VERY IMPORTANT CONCEPT
Go interfaces are implemented implicitly.
You NEVER write:
implements UserStore
like Java.
Instead:
If a struct has all required methods → it automatically satisfies the interface.
📦 Why is it accessible?
You asked:
Is it because they belong to package user?
✅ Partly yes.
Because:
package user
means:
routes.go
store.go
routes_test.go
all belong to SAME package namespace.
So:
they can access each other
they share types/functions
no import needed between them
🧠 The REAL reason the interface works
NOT because same package.
The REAL reason:
mockUserStore
implements the interface methods required by:
UserStore
That is the key.

*/

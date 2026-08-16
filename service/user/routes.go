package user

import (
	"log"
	"myproject/config"
	"myproject/service/auth"
	"myproject/types"
	"myproject/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

// Handler → struct that will handle user-related HTTP requests
type Handler struct {
	//store → field that holds a reference to the user store, which provides methods for interacting with user data (e.g., fetching, creating users)
	store types.UserStore
}

// NewHandler → constructor for the Handler struct
// Your handler does NOT depend on MySQL, it depends on the UserStore interface, which can be implemented by any storage mechanism (e.g., MySQL, PostgreSQL, in-memory store)
// You inject (pass) dependency from outside, which makes your code more flexible and easier to test
func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

// user package → handles user-related endpoints (/login, /register)
func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Register user-related routes
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

// handleLogin → processes login requests, validates credentials, and returns a token or session info
func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	//get json payload
	var payload types.LoginUserPayload
	// Parse the JSON payload from the request body into the RegisterUserPayload struct, returning an error response if parsing fails
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
	// Check if the user exists in the database by calling the GetUserByEmail method of the user store, returning an error response if the user is not found or if there is an error during the lookup

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	// Check if the provided password matches the stored hashed password using the ComparePasswords function from the auth package, returning an error response if the passwords do not match
	if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
		utils.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// If the user exists and the password matches, create a JWT token for the user by calling the CreateJWT function from the auth package, passing in the secret key and user ID. Return an error response if there is an error during token creation
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to create token")
		return
	}
// Return the generated JWT token in the response body as a JSON object with a "token" field, along with an HTTP status code of 200 (OK)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})

	
}

// handleRegister → processes registration requests, validates input, creates a new user, and returns success or error response
func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	//get json payload
	var payload types.RegisterUserPayload
	// Parse the JSON payload from the request body into the RegisterUserPayload struct, returning an error response if parsing fails
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
	//check if user exists

	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusConflict, "user already exists")
		return
	}
	// If user doesn't exist, hash the password before storing it
	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}
	//If doesn't exist, create user and return success response
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword, // In a real application, make sure to hash the password before storing it
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "failed to create user")
		return
	}
	log.Println("user created successfully:", payload.Email)
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user created successfully"})
}

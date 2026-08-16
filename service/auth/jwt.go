package auth

import (
	"context"
	"fmt"
	"myproject/config"
	"myproject/types"
	"myproject/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"

// CreateJWT creates a new JWT token with the given secret and user ID
// "I have already logged in. Here is my signed token. This token proves to the server which user I am."
// We use JWT (JSON Web Token) mainly to keep the user logged in and identify who is making requests without asking them to send their email and password every time.
func CreateJWT(secret []byte, userID int) (string, error) {

	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds) // Token expires in 24 hours
	// Create a new JWT token with the specified signing method and claims, including the user ID and expiration time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiresAt": time.Now().Add(expiration).Unix(), // Token expires in 24 hours
	})
	// Sign the token with the provided secret key and return the signed token string along with any error that occurs during signing
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// WithJWTAuth is a middleware function that checks for a valid JWT token in the request header and sets the user ID in the request context if the token is valid
func WithJWTAuth(handlerfunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	// Return a new http.HandlerFunc that wraps the provided handlerfunc with JWT authentication logic
	return func(w http.ResponseWriter, r *http.Request) {

		// get the token from the request header
		// The token is expected to be in the "Authorization" header, prefixed with "Bearer "
		tokenString, err := getTokenFromRequest(r)
		// gettokenfromrequest returns "bearer", later tirimpace remove bearer and put the authenticated toke like "ehyiturjgh"which is the actual jwt"
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, err.Error())
			permissionDenied(w)
			return
		}
		// validate the token
		// Validate the extracted token string using the validateToken function, which checks the token's signature and expiration. If the token is invalid or expired, return a 401 Unauthorized error and deny permission
		token, err := validateToken(tokenString)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "invalid token")
			permissionDenied(w)
			return
		}
		// Check if the token is valid. If not, return a 401 Unauthorized error and deny permission. If the token is valid, extract the user ID from the token claims, convert it to an integer, and set it in the request context for use in subsequent handlers
		if !token.Valid {
			utils.WriteError(w, http.StatusUnauthorized, "invalid token")
			permissionDenied(w)
			return
		}

		// if is , we need to fetch the userid from db
		// Extract the user ID from the token claims, which is stored as a string. Convert the string to an integer and handle any conversion errors. If successful, set the user ID in the request context for use in subsequent handlers
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)
		// conert "7" to int 7
		// Convert the user ID string to an integer. If the conversion fails, return a 401 Unauthorized error and deny permission
		userID, err := strconv.Atoi(str)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "failed to get userID from token")
			permissionDenied(w)
			return
		}
		// set context "userID" to the userid
		// Create a new context with the user ID and associate it with the request. This allows subsequent handlers to access the user ID from the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, userID)
		r = r.WithContext(ctx)

		handlerfunc(w, r)
		// handlerfunc calls the handleCheckout function, which processes the checkout request. The user ID is now available in the request context for use in the checkout process
		// get the string userid like"7" to int 7.
	}

	// so now the server knows the request belongs to userid 7.

}

// getTokenFromRequest extracts the JWT token from the request header, returning an error if the token is missing or malformed
func getTokenFromRequest(r *http.Request) (string, error) {
	// get the token from the request header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("missing token")
	}
	// Check if the token string starts with "Bearer " and remove the prefix if present, returning the trimmed token string
	if strings.HasPrefix(tokenString, "Bearer ") {
		// return the token string without the "Bearer " prefix, trimming any leading or trailing whitespace. This allows the token to be used for validation and authentication
		return strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer ")), nil
	}

	return tokenString, nil
}

// validateToken validates the JWT token using the provided secret, returning the parsed token and any error that occurs during validation
func validateToken(t string) (*jwt.Token, error) {
	// Validate the token using the jwt.Parse function, which checks the token's signature and expiration. The function returns the parsed token and any error that occurs during validation. If the token is valid, it can be used to extract claims such as the user ID
	return jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method used in the token matches the expected signing method (HMAC). If the signing method is unexpected, return an error indicating the mismatch. If the signing method is valid, return the secret key for signature verification
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key used for signing the token, which is retrieved from the configuration. This key is used to verify the token's signature and ensure its authenticity
		return []byte(config.Envs.JWTSecret), nil
	})
}

// permissionDenied writes a 403 Forbidden response with a "permission denied" message
func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, "permission denied")
}

// GetUserIDFromContext retrieves the user ID from the request context, returning an error if the user ID is not found
func GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return 0, fmt.Errorf("userID not found in context")
	}
	return userID, nil
}

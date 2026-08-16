package auth

import "golang.org/x/crypto/bcrypt"
// HashPassword → takes a plaintext password as input and returns a hashed version of the password using bcrypt, along with any error that occurs during the hashing process
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(hash), nil
}
// ComparePasswords() securely checks whether the password entered by the user matches the bcrypt-hashed password stored in the database, returning true if they match and false otherwise.
func ComparePasswords(hashed string, plain []byte ) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)
	return err == nil
	
}

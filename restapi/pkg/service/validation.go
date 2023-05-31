package service

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// Check credentials of user. Return true, nil if authentication is OK
func (s *service) Authentication(nickname, password string, userId uint16) (bool, error) {

	// Check if nickname exist
	isUserExist, err := s.ifUserExist(nickname)
	if err != nil {
		return false, fmt.Errorf("internal server error. Cannot confirm nickname exist. Error - %s", err.Error())
	}
	if isUserExist == false {
		return false, fmt.Errorf("user does not exist. Error - %s", err.Error())
	}

	// get User's hashed password from Database
	hashPassword, idDB, passErr := s.getPassword(nickname)
	if passErr != nil {
		return false, fmt.Errorf("internal server error - %s", passErr)
	}

	//check if provided user's id the same as in database for appropriate nickname
	if idDB != userId {
		return false, fmt.Errorf("incorrect id provided.")
	}

	// check if provided password is valid for  appropriate user ID and nickname
	valid := s.ComparePassword(hashPassword, password)
	if valid != nil {
		return false, fmt.Errorf("incorrect password.")
	}
	return true, nil
}

// check if user already exist in the database
func (s *service) ifUserExist(nickname string) (bool, error) {

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE nickname = $1)`

	var exist bool
	err := s.store.GetDB().DB.QueryRow(query, nickname).Scan(&exist)

	if err != nil {
		return false, err
	}
	return exist, nil
}

// get hashed password from the database for next comparison
func (s *service) getPassword(nickname string) (string, uint16, error) {

	query := `SELECT password, id FROM users WHERE nickname = $1`

	var password string
	var userId uint16

	err := s.store.GetDB().DB.QueryRow(query, nickname).Scan(&password, &userId)

	if err != nil {
		return "", 0, err
	}

	return password, userId, nil
}

// hashing of provided password during sing up
func (s *service) PasswordHashing(password string) (string, error) {

	isPasswordValid := s.PasswordValidation(password)
	if isPasswordValid != true {
		return "", errors.New(" the password does not meet the requirements")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// compare hashed and provided password
func (s *service) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// validate provided password.
// password must include at least 1 special symbol, 1 number and must be not shorten than 10 symbols
func (s *service) PasswordValidation(password string) bool {
	// minimum 10 characters
	if len(password) < 10 {
		return false
	}
	// at least one figure
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	if !hasNumber {
		return false
	}
	// at least one special character
	hasSpecialChar := regexp.MustCompile(`[^a-zA-Z0-9\s]`).MatchString(password)
	if !hasSpecialChar {
		return false
	}
	return true
}

package usecases

import (
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// Check credentials of user. Return true, nil if authentication is OK
func (repo *DBUserRepository) Authentication(nickname, password string, userid uint16) (bool, error) {

	// Check if nickname exist
	isUserExist, err := repo.ifUserExist(nickname)
	if err != nil {
		return false, fmt.Errorf("internal server error. Cannot confirm nickname exist. Error - $s", err)
	}
	if isUserExist == false {
		return false, fmt.Errorf("user does not exist. Error - %s", err)
	}

	// get User's hashed password from Database
	hashPassword, idDB, passErr := repo.getPassword(nickname)
	if passErr != nil {
		return false, fmt.Errorf("internal server error - %s", passErr)
	}

	//check if provided user's id the same as in database for appropriate nickname
	if idDB != userid {
		return false, fmt.Errorf("incorrect id provided.")
	}

	// check if provided password is valid for  appropriate user ID and nickname
	valid := repo.ComparePassword(hashPassword, password)
	if valid != nil {
		return false, fmt.Errorf("incorrect password. Error - %s", valid)
	}
	return true, nil
}

// check if user already exist in the database
func (repo *DBUserRepository) ifUserExist(nickname string) (bool, error) {
	var exist bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE nickname = $1)`
	err := repo.DB.QueryRow(query, nickname).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

// get hashed password from the database for next comparison
func (repo *DBUserRepository) getPassword(nickname string) (string, uint16, error) {
	var password string
	var userid uint16
	query1 := `select password from users where nickname = $1`
	query2 := `select id from users where nickname =$1`

	errPass := repo.DB.QueryRow(query1, nickname).Scan(&password)
	errID := repo.DB.QueryRow(query2, nickname).Scan(&userid)
	if errPass != nil {
		return "", 0, errPass
	}
	if errID != nil {
		return "", 0, errID
	}
	return password, userid, nil
}

// hashing of provided password during sing up
func (repo *DBUserRepository) PasswordHashing(password string) (string, error) {

	isPasswordValid := repo.PasswordValidation(password)
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
func (repo *DBUserRepository) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// validate provided password.
// password must include at least 1 special symbol, 1 number and must be not shorten than 10 symbols
func (repo *DBUserRepository) PasswordValidation(password string) bool {
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

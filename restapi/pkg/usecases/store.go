package usecases

import (
	"math"
	"os"

	"github.com/jmoiron/sqlx"

	log "github.com/sirupsen/logrus"

	"github.com/OlehEngineer/goTasks/goTasks/restapi/pkg/domian"
)

type DBUserRepository struct {
	DB *sqlx.DB
}

// interface which collect user processing methods
type Store interface {
	GetUsersPage(page, limit int) ([]domian.ApiResponse, error)
	GetPageQty(limit int) (int, error)
	GetUser(userid uint16) (domian.ApiResponse, error)
	CreateUser(newUser domian.UserSignUp, password string) (domian.ApiResponse, error)
	DeleteUser(userid uint16) error
	UpdateUser(upUser domian.UpdateUser) (domian.ApiResponse, error)
	GetDB() *sqlx.DB
}

func New() *DBUserRepository {
	//connect to the Postgres database and check Ping
	conn, err := DatabaseConnect()
	if err != nil {
		log.Fatalf("%s - %v", os.Getenv("DATABASECONNECTFAIL"), err)
	}
	return &DBUserRepository{DB: conn}
}

// get users from the database. Pagination implemented.
func (repo *DBUserRepository) GetUsersPage(page, limit int) ([]domian.ApiResponse, error) {

	query := `SELECT id, nickname, name, lastname, email, status, created_at, updated_at, deleted_at, likes FROM users ORDER BY id LIMIT $1 OFFSET $2`

	var offset int
	// calculate offset.
	switch {
	case page == 1:
		offset = 0
	case page > 1:
		offset = (page - 1) * limit
	}

	var users []domian.ApiResponse
	err := repo.DB.Select(&users, query, limit, offset)
	if err != nil {
		log.Errorf("GetUsersPage. Select error - %v", err)
		return users, err
	}

	return users, nil
}

// get user by id. Basic Authentication implemented
func (repo *DBUserRepository) GetUser(userid uint16) (domian.ApiResponse, error) {

	query := `SELECT id, nickname, name, lastname, email, status, created_at, updated_at, deleted_at, likes FROM users WHERE id = $1`

	user := domian.ApiResponse{}
	err := repo.DB.Get(&user, query, userid)
	if err != nil {
		log.Errorf("GetUser Get error - %s", err)
		return user, err
	}

	return user, nil
}

// create new user. Basic Authentication implemented
func (repo *DBUserRepository) CreateUser(newUser domian.UserSignUp, password string) (domian.ApiResponse, error) {

	var id int16
	// created := time.Now().UTC().String()
	query := `INSERT INTO users (nickname, name, lastname, email, password, status, created_at, updated_at, deleted_at, likes) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, $7) RETURNING id `

	user := domian.ApiResponse{}
	err := repo.DB.QueryRow(query, newUser.NickName, newUser.Name, newUser.LastName, "noEmail", password, "active", 0).Scan(&id)
	if err != nil {
		log.Errorf("PostUser QueryRow error - %s", err)
		return user, err
	}
	user, err = repo.GetUser(uint16(id))
	if err != nil {
		return user, err
	}
	return user, nil
}

// delete user by id. Basic Authentication implemented
func (repo *DBUserRepository) DeleteUser(userid uint16) error {

	query := `DELETE FROM users WHERE id =$1`

	_, err := repo.DB.Exec(query, userid)
	if err != nil {
		log.Errorf("DeleteUser Exec error - %s", err)
		return err
	}
	return nil
}

// update user by id. Basic Authentication implemented
func (repo *DBUserRepository) UpdateUser(upUser domian.UpdateUser) (domian.ApiResponse, error) {

	query := `UPDATE users SET  name=$1, lastname=$2, email=$3, updated_at=CURRENT_TIMESTAMP, likes=$4 WHERE id=$5`

	user := domian.ApiResponse{}
	_, err := repo.DB.Exec(query, upUser.Name, upUser.LastName, upUser.EmailAddress, upUser.Likes, upUser.UserID)
	if err != nil {
		log.Errorf("PutUser Exec error - %s", err)
		return user, err
	}
	user, err = repo.GetUser(uint16(upUser.UserID))
	if err != nil {

		return user, err
	}
	return user, nil
}

// return pages quantity based on limit per page
func (repo *DBUserRepository) GetPageQty(limit int) (int, error) {

	query := `SELECT count(id) FROM users`

	var totalRows int
	err := repo.DB.Get(&totalRows, query)
	if err != nil {
		log.Errorf("GetRowQtyDB. cannot get rows q-ty from the database. error - %s", err)
		return 0, err
	}
	// Calculate the total number of pages
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	return totalPages, nil
}
func (repo *DBUserRepository) GetDB() *sqlx.DB {
	return repo.DB
}

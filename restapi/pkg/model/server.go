package model

import (
	"math"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

func GetUsersPage(conn *sqlx.DB, page, limit int) ([]ApiResponse, error) {
	var users []ApiResponse
	var offset int
	query := `select id, nickname, name, lastname, email, status, created_at, updated_at, delete_at, likes from users ORDER BY id LIMIT $1 OFFSET $2`

	// calculate offset.
	switch {
	case page == 1:
		offset = 0
	case page > 1:
		offset = (page - 1) * limit
	}

	err := conn.Select(&users, query, limit, offset)
	if err != nil {
		log.Errorf("GetUsersPage. Select error - %v", err)
		return users, err
	}

	return users, nil
}

func GetUser(conn *sqlx.DB, userid uint16) (ApiResponse, error) {
	user := ApiResponse{}
	query := `select id, nickname, name, lastname, email, status, created_at, updated_at, delete_at, likes from users where id = $1`

	err := conn.Get(&user, query, userid)
	if err != nil {
		log.Errorf("GetUser Get error - %s", err)
		return user, err
	}

	return user, nil
}

func PostUser(conn *sqlx.DB, nickname, name, lastname, password string) (ApiResponse, error) {
	newUser := ApiResponse{}
	var id int16
	created := time.Now().UTC().String()
	query := `insert into users (nickname, name, lastname, email, password, status, created_at, updated_at, delete_at, likes) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id `

	err := conn.QueryRow(query, nickname, name, lastname, "noEmail", password, "active", created, "", "", 0).Scan(&id)
	if err != nil {
		log.Errorf("PostUser QueryRow error - %s", err)
		return newUser, err
	}
	newUser, err = GetUser(conn, uint16(id))
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func DeleteUser(conn *sqlx.DB, userid uint16) error {

	query := `delete from users where id =$1`

	_, err := conn.Exec(query, userid)
	if err != nil {
		log.Errorf("DeleteUser Exec error - %s", err)
		return err
	}
	return nil
}
func PutUser(conn *sqlx.DB, upUser UpdateUser) (ApiResponse, error) {

	updatedUser := ApiResponse{}
	updated := time.Now().UTC().String()
	query := `UPDATE users SET  name=$1, lastname=$2, email=$3, updated_at=$4, likes=$5 WHERE id=$6`

	_, err := conn.Exec(query, upUser.Name, upUser.LastName, upUser.EmailAddress, updated, upUser.Likes, upUser.UserID)
	if err != nil {
		log.Errorf("PutUser Exec error - %s", err)
		return updatedUser, err
	}
	updatedUser, err = GetUser(conn, uint16(upUser.UserID))
	if err != nil {

		return updatedUser, err
	}
	return updatedUser, nil
}
func GetPageQty(conn *sqlx.DB, limit int) (int, error) {

	var totalRows int
	query := `select count(id) from users`

	err := conn.Get(&totalRows, query)
	if err != nil {
		log.Errorf("GetRowQtyDB. cannot get rows q-ty from the database. error - %s", err)
		return 0, err
	}
	// Calculate the total number of pages
	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))

	return totalPages, nil
}

package usecases

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// Connect to the database
func DatabaseConnect() (*sqlx.DB, error) {
	// get configuration"
	dbHost := os.Getenv("HOST")
	dbPort := os.Getenv("PORT")
	dbName := os.Getenv("DBNAME")
	dbUser := os.Getenv("USER")
	dbPassword := os.Getenv("PASSWORD")
	dbSSLmode := os.Getenv("SSLMODE")

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPassword, dbSSLmode)

	//start connection to database
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("cannot connect to the database. Error - %s", err)
		return nil, err
	}
	log.Info("database connected successfully")

	//ping database
	err = db.Ping()
	if err != nil {
		log.Fatal("ping nok. Lost connection")
		db.Close()
		return nil, err
	}
	return db, nil
}

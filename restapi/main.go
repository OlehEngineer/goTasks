package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/OlehEngineer/goTasks/goTasks/restapi/pkg/logger"
	"github.com/OlehEngineer/goTasks/goTasks/restapi/pkg/model"
)

func main() {

	// Load environment variables
	if confErr := godotenv.Load(); confErr != nil {
		log.Fatalf("Cannot read .ENV file. Error - %v", confErr)
	}

	//Starting logging
	logger.StartLogging(os.Getenv("LOGGERLEVEL"))
	log.Info(os.Getenv("LOGSTART"))

	//connect to the Postgres database and check Ping
	conn, err := model.DatabaseConnect()
	if err != nil {
		log.Fatalf("%s - %v", os.Getenv("DATABASECONNECTFAIL"), err)
	}
	defer conn.Close()

	//Start new ECHO
	e := echo.New()
	model.RegisterRouters(e)
	// Inject the database pool into the Echo instance
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("conn", conn)
			return next(c)
		}
	})
	//start server
	e.Logger.Fatal(e.Start(":8080"))

}

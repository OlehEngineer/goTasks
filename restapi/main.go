package main

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/OlehEngineer/goTasks/goTasks/restapi/pkg/controllers"
	"github.com/OlehEngineer/goTasks/goTasks/restapi/pkg/service"
	"github.com/OlehEngineer/goTasks/goTasks/restapi/pkg/usecases"
)

func main() {
	// Load environment variables
	if confErr := godotenv.Load(); confErr != nil {
		log.Fatalf("Cannot read .ENV file. Error - %v", confErr)
	}
	//Starting logging
	usecases.StartLogging(os.Getenv("LOGGERLEVEL"))
	log.Info(os.Getenv("LOGSTART"))

	//create usecase layer
	useCaseLayer := usecases.New()

	//create service layer
	serviceLayer := service.New(useCaseLayer)

	//create router layer
	api := controllers.New(serviceLayer)

	api.Run()

}

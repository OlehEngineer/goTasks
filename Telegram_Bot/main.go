package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	//create Log file for errors
	LogFile, errLog := os.OpenFile("LOG.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	defer LogFile.Close()
	//LOGRUS setting
	multiOutput := io.MultiWriter(os.Stdout, LogFile) //set logging into standard output and into a file
	log.SetOutput(multiOutput)
	log.SetLevel(log.TraceLevel) //Set log level
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if errLog != nil {
		log.Fatal("problem with errors log file")
	}
	//real .ENV file with access TOKEN and settings
	envErr := godotenv.Load("local.env")
	if envErr != nil {
		log.Fatalf(".ENV file is missing. Error > %s < occurs.", envErr)
	}

	botToken := os.Getenv("TOKEN") //please insert your API TOKEN here
	botApi := "https://api.telegram.org/bot"
	botURL := botApi + botToken
	offset := 0
	//↓ infinity cycle for updates checking ↓
	for {
		updates, err := getUpdates(botURL, offset)
		if err != nil {
			log.Errorf("Something went wrong: ", err.Error())
		}
		//iteration throught each element of updates
		//method sendMassage is using
		for _, update := range updates {
			err = respond(botURL, update)
			if err != nil {
				log.Error(err)
			}
			//new updates id. Increase for each new update
			offset = update.UpdateId + 1
		}
		// write logging parameters
		for _, upd := range updates {
			log.Infof("Update ID - %v, Chat ID - %v, Massage - «%s»\n", upd.UpdateId, upd.Message.Chat.ChatId, upd.Message.Text)
		}

	}
}

// ask Updates
func getUpdates(botURL string, offset int) ([]Update, error) {
	//reques to the Bot with method "getUpdates"
	resp, err := http.Get(botURL + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	// close Bot respond body in the end of this function
	defer resp.Body.Close()
	//transtale bytes format to readble format
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//parsing of json Bot's reply
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

// respond to User
func respond(botUrl string, update Update) error {
	var botMessage BotMessage

	botMessage.ChatId = update.Message.Chat.ChatId

	botMessage.Text = update.Message.Text

	defaultAnswer := fmt.Sprintf("there is not such command: >%s<. Please input /help to get a list with available commands", update.Message.Text)
	//check user input
	switch botMessage.Text {
	case "/about":
		botMessage.Text = os.Getenv("MYINFO") //"I ma 33 year old engineer who wants to became a programmer"
	case "/links":
		botMessage.Text = fmt.Sprintf("My Github account - %s\nMy personal email address is - %s", os.Getenv("MYGITHUB"), os.Getenv("MYMAIL"))
	case "/start":
		botMessage.Text = "possible commands: /about; /links; /start; /help"
	case "/help":
		botMessage.Text = "possible commands: /about; /links; /start"
	default:
		botMessage.Text = defaultAnswer
	}
	//convert respond messagge in bytes format
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	//send answer to user using method "sendMessage" in bytes format
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

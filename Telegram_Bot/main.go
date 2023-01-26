package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load("local.env")
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TOKEN") //please insert your API TOKEN here
	botApi := "https://api.telegram.org/bot"
	botURL := botApi + botToken
	offset := 0
	//↓ infinity cycle for updates checking ↓
	for {
		updates, err := getUpdates(botURL, offset)
		if err != nil {
			log.Println("Something went wrong: ", err.Error())
		}
		//iteration throught each element of updates
		//method sendMassage is using
		for _, update := range updates {
			err = respond(botURL, update)
			if err != nil {
				log.Println(err)
			}
			//new updates id. Increase for each new update
			offset = update.UpdateId + 1
		}
		fmt.Println(updates)
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
	//send the same message to user which Bot gets
	botMessage.Text = update.Message.Text
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
		botMessage.Text = "there is not such command. Please input /help to get a list with available commands"
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

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// respond to User
func respond(botUrl string, update Update) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = update.Message.Text
	defaultAnswer := fmt.Sprintf("there is not such command: >%s<. Please input /help to get a list with available commands", update.Message.Text)
	//check user input
	switch botMessage.Text {
	case "/about":
		botMessage.Text = CurrentConfiguration.Myinfo //get MYINFO from .ENV file
	case "/links":
		botMessage.Text = fmt.Sprintf("My Github account - %s\nMy personal email address is - %s", CurrentConfiguration.MuGitHub, CurrentConfiguration.Myemail) //get MYGIHUB and MYMAIL infromation from .ENV file
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

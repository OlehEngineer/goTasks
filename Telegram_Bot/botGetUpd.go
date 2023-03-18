package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

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

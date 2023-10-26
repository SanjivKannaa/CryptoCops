package helpers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Link struct {
	Url string `json:"url"`
}

func Write(links []Link, filePath string) error {
	dataBytes, err := json.Marshal(links)
	if err != nil {
		return errors.New("Error marshalling the json file: " + err.Error())
	}
	err = ioutil.WriteFile(filePath, dataBytes, 0644)
	if err != nil {
		return errors.New("Error writing to the url queue file: " + err.Error())
	}
	return nil
}

func Read(filePath string) ([]Link, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("Error opening the url queue file: " + err.Error())
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.New("Error reading the url queue file: " + err.Error())
	}
	var links []Link
	err = json.Unmarshal(byteValue, &links)
	if err != nil {
		return nil, errors.New("Error unmarshalling the json file: " + err.Error())
	}
	return links, err
}

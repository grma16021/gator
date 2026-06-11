package config

import (
	"encoding/json"
	"os"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

var conf Config

func Read() (Config, error) {

	homeDirect, err := os.UserHomeDir()
	if err != nil {
		println("error getting home direct", homeDirect, err)
		return Config{}, err
	}

	var filepath = homeDirect + configFileName
	fileContent, err := os.ReadFile(homeDirect + configFileName)
	if err != nil {
		println("Error reading", filepath, err)
		return Config{}, err
	}
	err = json.Unmarshal(fileContent, &conf)
	if err != nil {
		println("Error unmarshaling", err)
		return Config{}, err
	}

	return conf, nil

}
func (c Config) SetUser(name string) {
	c.Current_user_name = name
	homeDirect, err := os.UserHomeDir()
	var filepath = homeDirect + configFileName
	b, err := json.Marshal(c)
	if err != nil {
		println("error marshaling")
	}

	err = os.WriteFile(filepath, b, 0644)
	if err != nil {
		println("error writing to file")
	}
}

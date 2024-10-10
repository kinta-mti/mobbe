package config

import (
	"encoding/json"
	"log"
	"os"
)

type EnvironmentInfo struct {
	Name        string `json:"name"`
	OrderPrefix string `json:"orderPrefix"`
}

type DBConnInfo struct {
	Name string `json:"name"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

type ServerInfo struct {
	Port string `json:"port"`
}

type YpgInfo struct {
	ApiKey     string `json:"apiKey"`
	SecretKey  string `json:"SecretKey"`
	Apimkey    string `json:"apimKey"`
	ApimSecret string `json:"apimSecret"`
	Uri        string `json:"uri"`
	Path       Path   `json:"path"`
}

type Path struct {
	AccesToken string `json:"accesToken"`
	Inquiries  string `json:"inquiries"`
}

type Configuration struct {
	Environment EnvironmentInfo `json:"environment"`
	Server      ServerInfo      `json:"server"`
	Database    DBConnInfo      `json:"database"`
	Ypg         YpgInfo         `json:"ypg"`
}

func Load(path string) Configuration {
	file, _ := os.Open(path)
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Println("error:", err)
	}
	log.Println("Server Info:" + configuration.Server.Port)
	log.Println("DB Connection Info:" + configuration.Database.Name + ";" + configuration.Database.User + ";" + configuration.Database.Pass)
	return configuration
}

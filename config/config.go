package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("error loading config file")
	}
}

func DBConnString() string {
	host := os.Getenv("DBHOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASSWORD")
	dbname := os.Getenv("DBNAME")
	sslmode := os.Getenv("DBSSLMODE")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
}

func ServerConfig() string {
	//host := os.Getenv("SERVERHOST")
	port := os.Getenv("SERVERPORT")

	return port
}

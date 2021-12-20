package database

import (
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/matsuokashuhei/go-todo/config"
)

func ConnectDB() {
	var err error
	p := config.Config("POSTGRES_PORT")
	port, err := strconv.ParseUint(p, 10, 32)

	configData := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Config("POSTGRES_HOST"),
		port,
		config.Config("POSTGRES_USER"),
		config.Config("POSTGRES_PASSWORD"),
		config.Config("POSTGRES_NAME"),
	)

	DB, err = gorm.Open("postgres", configData)

	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
}

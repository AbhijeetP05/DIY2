package main

import (
	"fmt"
	"os"
)

func main() {
	a := App{}
	//fmt.Println("host is " + os.Getenv("APP_DB_HOST"))
	a.Initialize(os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_PORT"),
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))
	fmt.Println("Database connected and starting server at 8010...")
	a.Run(os.Getenv("APP_DB_HOST"), "8010")
}

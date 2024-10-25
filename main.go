package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

var DisplayName string = "Unset"
var ShortName string = "unset"
var Version string = "?.?.?"
var Commit string = "???????"

func main() {
	err := godotenv.Load()
	if err == nil {
		fmt.Println("Environment variables from .env loaded")
	}

	log.Println(DisplayName + " version " + Version + ", build " + Commit)

}

package config_test

import (
	"fmt"
	"cli_sederhana2/config"
	"testing"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load("cli_sederhana2/.env")
	if err != nil {
		fmt.Println("env not found, using system env")
	}
}

func TestConnection(t *testing.T) {
	Init()
	config.OpenDB()
}

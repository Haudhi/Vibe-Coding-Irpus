package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	fmt.Println("=== Environment Variables ===")
	fmt.Printf("DB_HOST: '%s'\n", os.Getenv("DB_HOST"))
	fmt.Printf("DB_PORT: '%s'\n", os.Getenv("DB_PORT"))
	fmt.Printf("DB_NAME: '%s'\n", os.Getenv("DB_NAME"))
	fmt.Printf("DB_USER: '%s'\n", os.Getenv("DB_USER"))
	fmt.Printf("DB_PASSWORD: '%s'\n", os.Getenv("DB_PASSWORD"))

	// Set up viper like in config
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	fmt.Println("\n=== After Viper Setup ===")
	fmt.Printf("Database Host: '%s'\n", viper.GetString("database.host"))
	fmt.Printf("Database Port: %d\n", viper.GetInt("database.port"))
	fmt.Printf("Database Name: '%s'\n", viper.GetString("database.database"))
	fmt.Printf("Database User: '%s'\n", viper.GetString("database.user"))
	fmt.Printf("Database Password: '%s'\n", viper.GetString("database.password"))
}
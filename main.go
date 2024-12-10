package main

import (
	"fmt"
	"log"
	"os"

	"travel-planner/cmd"

	"github.com/spf13/viper"
)

func main() {
	// Load configuration
	viper.SetConfigFile("configs/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Execute the root command and directly handle any errors
	if err := cmd.Execute(); err != nil {
		// If an error occurs, print the error and exit with a non-zero status
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}

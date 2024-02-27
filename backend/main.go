package main

import (
	"fmt"

	// colorboxd "github.com/dsantos747/letterboxd_hue_sort/backend/colorboxd"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Could not load environment variables from .env file: %v\n", err)
		return
	}

}

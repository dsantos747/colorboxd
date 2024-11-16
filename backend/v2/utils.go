package colorboxd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv attempts to load an env var "ENVIRONMENT". If successful, no further action.
// If not successful, load all envs with godotenv instead
func LoadEnv() error {
	if os.Getenv("ENVIRONMENT") == "" {
		err := godotenv.Load("../.env")
		if err != nil {
			fmt.Printf("Could not load environment variables from .env file: %v\n", err)
			return err
		}
	}
	return nil
}

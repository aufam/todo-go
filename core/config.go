package core

import (
	"fmt"
	"os"
)

func GetVersion() (string, error) {
	uri, err := LoadEnv("API_VERSION")
	if err == nil {
		return "v1", nil
	}

	return uri, err
}

func GetMongoDBURI() (string, error) {
	uri, err := LoadEnv("MONGODB_URI")
	if err == nil {
		return uri, nil
	}

	username, err := LoadEnv("MONGO_INITDB_ROOT_USERNAME")
	if err != nil {
		return "", err
	}

	password, err := LoadEnv("MONGO_INITDB_ROOT_PASSWORD")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("mongodb://%s:%s@mongodb:27017", username, password), nil
}

func LoadEnv(e string) (string, error) {
	v := os.Getenv(e)
	if v == "" {
		return "", fmt.Errorf("Environment variable `%s` is not provided", e)
	}

	return v, nil
}

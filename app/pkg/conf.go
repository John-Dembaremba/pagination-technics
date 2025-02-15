package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Env struct {
	ProjectVersion   string `mapstructure:"PROJECT_VERSION"`
	ServerPort       string `mapstructure:"SERVER_PORT"`
	POSTGRES_VERSION string `mapstructure:"POSTGRES_VERSION"`
	POSTGRES_DB      string `mapstructure:"POSTGRES_DB"`
	POSTGRES_USER    string `mapstructure:"POSTGRES_USER"`
	POSTGRES_PSW     string `mapstructure:"POSTGRES_PSW"`
	POSTGRES_PORT    string `mapstructure:"POSTGRES_PORT"`
	POSTGRES_HOST    string `mapstructure:"POSTGRES_HOST"`
}

func NewEnv() Env {
	env := &Env{}

	envPath, err := findEnvFilePath()
	if err != nil {
		log.Fatal("Can't find the file .env in any parent directory : ", err)
	}

	viper.SetConfigFile(envPath)
	viper.SetConfigType("env")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)

	}

	return *env
}

func findEnvFilePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	for {
		log.Printf("Searching for .env file starting from: %s", dir)

		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf(".env file not found")
}

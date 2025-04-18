package pkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Env struct {
	ProjectVersion          string `mapstructure:"PROJECT_VERSION"`
	ServerPort              string `mapstructure:"SERVER_PORT"`
	POSTGRES_CONTAINER_NAME string `mapstructure:"POSTGRES_CONTAINER_NAME"`
	POSTGRES_VERSION        string `mapstructure:"POSTGRES_VERSION"`
	POSTGRES_DB             string `mapstructure:"POSTGRES_DB"`
	POSTGRES_USER           string `mapstructure:"POSTGRES_USER"`
	POSTGRES_PSW            string `mapstructure:"POSTGRES_PSW"`
	POSTGRES_PORT           string `mapstructure:"POSTGRES_PORT"`
	POSTGRES_HOST           string `mapstructure:"POSTGRES_HOST"`

	JAEGER_HOST    string `mapstructure:"JAEGER_HOST"`
	OTLP_HTTP_PORT int    `mapstructure:"OTLP_HTTP_PORT"`
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

func ReadFile(path string) (string, error) {
	// Get the absolute path of the schema file
	absPath, err := filepath.Abs(filepath.Join("pkg", path))
	if err != nil {
		return "", err
	}

	file, err := os.Open(absPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

package viper

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type getFromEnvFn func(string) string

// ConfigService is a Viper config service
type ConfigService struct {
	*viper.Viper
}

// New creates an initialized ConfigService
func New(getFromEnv getFromEnvFn) (ConfigService, error) {
	config := ConfigService{
		viper.New(),
	}

	configFullPath, err := getHomePath(getFromEnv("CONFIG_PATH"))
	if err != nil {
		return ConfigService{}, err
	}

	config.SetConfigFile(fmt.Sprintf("%s/%s.yaml", configFullPath, getFromEnv("GOENV")))
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		return ConfigService{}, fmt.Errorf("fatal error reading viper file: %s", err)
	}

	return config, nil
}

// GetString returns the value associated with the key as a string
func (config ConfigService) GetString(key string) string {
	return config.Viper.GetString(key)
}

// GetInt returns the value associated with the key as a int
func (config ConfigService) GetInt(key string) int {
	return config.Viper.GetInt(key)
}

// GetBool looks for the given key to return a bool
func (config ConfigService) GetBool(key string) bool {
	return config.Viper.GetBool(key)
}

func getHomePath(configPath string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	dir := usr.HomeDir

	if configPath == "~" {
		return dir, nil
	}

	if strings.HasPrefix(configPath, "~/") {
		return filepath.Join(dir, configPath[2:]), nil
	}

	return configPath, nil
}

package it_test

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func getConfigFilePath(fileName string) string {
	dir, _ := os.Getwd()
	return filepath.FromSlash(fmt.Sprintf("%s/testdata/%s", dir, fileName))
}

func setConfig(cfgFile string) error {
	viper.Reset()
	viper.SetConfigFile(cfgFile)
	err := viper.ReadInConfig()

	if err != nil {
		return fmt.Errorf("failed to load '%s' configuration file. Error %s", cfgFile, err)
	}

	return nil
}

func deleteConfig(cfgFile string) error {
	return os.Remove(cfgFile)
}

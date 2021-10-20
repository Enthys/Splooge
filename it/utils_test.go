package it_test

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"wildfire/pkg"
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

func initiateConfiguration(cfgFile string) error {
	_ = os.Remove(cfgFile)
	_ = setConfig(cfgFile)
	config := pkg.GetConfig()
	projectService := pkg.NewProjectService(config)
	projectService.UpdateOrCreate(&pkg.ProjectConfig{"foo", pkg.ProjectTypeGit, "github.com/url"})
	projectService.UpdateOrCreate(&pkg.ProjectConfig{"bar", pkg.ProjectTypeGit, "github.com/url"})
	projectService.UpdateOrCreate(&pkg.ProjectConfig{"zaz", pkg.ProjectTypeGit, "github.com/url"})
	err := config.SaveConfig()
	if err != nil {
		return fmt.Errorf("failed to save configuration")
	}
	config = pkg.GetConfig()
	if len(config.Projects) != 3 {
		return fmt.Errorf(
			"expected to have 3 projects found '%d'",
			len(config.Projects),
			)
	}

	return nil
}

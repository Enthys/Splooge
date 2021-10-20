package pkg

import (
	"github.com/spf13/viper"
)

type WildFireConfig struct {
	Projects map[string]*ProjectConfig `yaml:"projects"`
	Groups map[string]*GroupConfig     `yaml:"groups"`
}

func GetConfig() *WildFireConfig {
	var config WildFireConfig
	err := viper.Unmarshal(&config)

	if err != nil {
		return &WildFireConfig{
			Projects: make(map[string]*ProjectConfig),
			Groups: make(map[string]*GroupConfig),
		}
	}

	if len(config.Projects) == 0 {
		config.Projects = make(map[string]*ProjectConfig)
	}

	if len(config.Groups) == 0 {
		config.Groups = make(map[string]*GroupConfig)
	}

	return &config
}

func (config *WildFireConfig) SaveConfig() error {
	viper.Set("projects", config.Projects)
	viper.Set("groups", config.Groups)

	return viper.WriteConfig()
}

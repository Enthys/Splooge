package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"splooge/pkg"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "splooge",
	Short: "Application for mass update of repositories",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var addProjectCmd = &cobra.Command{
	Use: "add-project [name] [type] [url]",
	Short: "Add Project to the loaded configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Errorf("%s", "Invalid number of arguments provided")
		}

		config := pkg.GetConfig()
		err := config.AddProject(&pkg.Project{
			Name: args[0],
			Type: pkg.ProjectType(args[1]),
			URL:  pkg.ProjectPath(args[2]),
		})

		cobra.CheckErr(err)

		err = config.SaveConfig()
		cobra.CheckErr(err)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.splooge.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(addProjectCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		dir, err := os.Getwd()
		cobra.CheckErr(err)

		// Search config in current directory with name ".splooge" (without extension).
		cfgFile = filepath.FromSlash(dir + "/.splooge.yaml")
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Configuration file", cfgFile, "failed to load:", err)
	} else {
		fmt.Println("Using configuration file", cfgFile)
	}
}

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"wildfire/cmd/group"
	"wildfire/cmd/project"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "wildfire",
	Short: "Application for mass update of repositories",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wildfire.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(project.ProjectCmd)
	rootCmd.AddCommand(group.GroupCmd)
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

		// Search config in current directory with name ".wildire" (without extension).
		cfgFile = filepath.FromSlash(dir + "/.wildfire.yaml")
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

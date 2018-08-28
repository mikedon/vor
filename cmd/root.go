package commands

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trevor-atlas/vor/env"
	"github.com/trevor-atlas/vor/formatters"
	"github.com/trevor-atlas/vor/logger"
	"github.com/trevor-atlas/vor/system"
	"os"
	"path/filepath"
	"strings"
)

var (
	CONFIG_FILE string
)

var rootCmd = &cobra.Command{
	Use:   "vor",
	Short: "Vör – make Github and Jira easy",
	Long: `
                  ___          ___
      ___        /\  \        /\  \
     /\  \      /::\  \      /::\  \
     \:\  \    /:/\:\  \    /:/\:\__\
      \:\  \  /:/  \:\  \  /:/ /:/  /
  ___  \:\__\/:/__/ \:\__\/:/_/:/__/___
 /\  \ |:|  |\:\  \ /:/  /\:\/:::::/  /
 \:\  \|:|  | \:\  /:/  /  \::/~~/~~~~
  \:\__|:|__|  \:\/:/  /    \:\~~\
   \::::/__/    \::/  /      \:\__\
    ~~~~         \/__/        \/__/

 This program comes with ABSOLUTELY NO WARRANTY; This is free software, and you are welcome to redistribute it.
 Vör – A fast and flexible commandline tool for working with Github and Jira`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		env.Init(&env.DefaultLoader{})
		formatters.Init(&formatters.DefaultStringFormatter{})
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	//  If Persistent flags are defined here they are global
	rootCmd.PersistentFlags().StringVar(&CONFIG_FILE, "config", "", "config file (default is $HOME/.vor, or the current directory)")
	viper.SetDefault("devmode", false)
	viper.SetDefault("git.branchtemplate", "{jira-issue-number}/{jira-issue-type}/{jira-issue-title}")
	viper.SetDefault("git.path", "/usr/local/bin/git")
	viper.SetDefault("git.pull-request-base", "master")
}

func initConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName(".vor")

	if CONFIG_FILE != "" {
		viper.SetConfigFile(CONFIG_FILE) // Use config file from the flag.
		return
	}

	home, homeErr := homedir.Dir()
	if homeErr != nil {
		fmt.Println(homeErr)
		system.Exit("vor encountered an error attempting to read from the filesystem")
	}
	configPath, walkErr := walkUpFS(filepath.Base(home), )
	fmt.Println("configPath: "+configPath)
	if walkErr != nil {
		fmt.Println(walkErr)
		system.Exit("vor encountered an error attempting to read from the filesystem")
	}

	viper.AddConfigPath(home)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		color.Red("Vor could not find a local config file, this can cause problems and is not recommended\n")
		fmt.Println(err)
	}
}

func walkUpFS(cutoff string) (match string, err error) {
	// get the directory we started in
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// split it into directories
	path := strings.Split(currentDir, "/")
	cutoff = filepath.Clean(cutoff)
	for i := len(path) - 1; i >= 0; i-- {
		here := strings.Join(path[0:i], "/")
		if cutoff == here || i == 0 { // if we are at the cutoff directory, stop
			return "", nil
		} else { // otherwise, walk the path from i
			logger.Debug("walking " + here)
			matches, err := filepath.Glob(here + "/.vor.*")
			if err != nil{
				logger.Debug("walk error [%v]\n", err)
				return "", err
			}
			if len(matches) > 0 {
				logger.Debug("walker done!")
				return matches[0], nil
			}
		}
	}
	logger.Debug("walker done!")
	return "", nil
}

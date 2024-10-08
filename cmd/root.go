/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	utils "met/utils"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:               "my-ether-tool",
	Short:             "my ether tool",
	Long:              `evm based blockchain client CLI tool`,
	PersistentPreRunE: rootPreRun,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.my-ether-tool.yaml)")
	RootCmd.PersistentFlags().String("loglevel", "", "log level: trace debug info warn error fatal (loglevel has high priority than verbose)")
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose message (debug loglevel)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func rootPreRun(cmd *cobra.Command, args []string) error {
	loglevel := cmd.Flag("loglevel").Value.String()
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		panic(err)
	}

	var levelStr string

	// loglevel 优先级高于verbose
	if loglevel != "" {
		levelStr = loglevel
	} else {
		if verbose {
			levelStr = "debug"
		} else {
			levelStr = "info"
		}
	}

	err = utils.SetLogger(levelStr)
	if err != nil {
		return err
	}

	return nil
}

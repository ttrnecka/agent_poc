package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initConfig() error {
	log.Printf("Initializing configuration for %s version %s", NAME, VERSION)
	log.Print("Checking environment variables and flags")
	log.Print("Checking if CLI_USER is set")
	err := viper.BindEnv("user", "CLI_USER")

	if err != nil {
		return fmt.Errorf("failed to bind CLI_USER: %w", err)
	}

	if viper.GetString("user") == "" {
		return fmt.Errorf("CLI_USER environment variable is required but not set")
	}
	log.Printf("Using user: %s", viper.GetString("user"))

	log.Print("Checking if CLI_PASSWORD is set")
	err = viper.BindEnv("password", "CLI_PASSWORD")

	if err != nil {
		return fmt.Errorf("failed to bind CLI_PASSWORD: %w", err)
	}

	if viper.GetString("password") == "" {
		return fmt.Errorf("CLI_PASSWORD environment variable is required but not set")
	}
	log.Printf("Password is set")

	err = checkFolder(viper.GetString("output_folder"))
	if err != nil {
		return fmt.Errorf("failed to check output folder: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to check working directory: %w", err)
	}
	log.Printf("Output folder is set to: %s", filepath.Join(wd, viper.GetString("output_folder")))

	return nil
}

var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s_%s", NAME, VERSION),
	Short: "Brocade Collector Plugin",
	Long:  `Brocade Collector Plugin`,
	Run: func(cmd *cobra.Command, args []string) {
		// We ignore default as we only user subcommands
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
	// SilenceErrors: true,
}

func Execute() {

	rootCmd.PersistentFlags().StringP("endpoint", "e", "", "host:port notation of the endpoint to connect to")
	rootCmd.MarkFlagRequired("endpoint")
	viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))

	rootCmd.PersistentFlags().StringP("output_folder", "o", "", "folder to store collected data")
	rootCmd.MarkFlagRequired("output_folder")
	viper.BindPFlag("output_folder", rootCmd.PersistentFlags().Lookup("output_folder"))

	collectCmd := NewCollectCmd()
	validateCmd := NewValidateCmd()

	rootCmd.AddCommand(collectCmd)
	rootCmd.AddCommand(validateCmd)

	if err := rootCmd.Execute(); err != nil {
		// no need as execute will print the error by default
		// fmt.Fprintln(os.Stderr, "Error:", err)
		if ec, ok := err.(*exitCodeError); ok {
			os.Exit(ec.Code)
		} else {
			os.Exit(1)
		}
	}
}

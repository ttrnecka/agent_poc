package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	logging "github.com/ttrnecka/agent_poc/logger"
	"golang.org/x/mod/semver"
)

var logger zerolog.Logger

var done chan bool

type Cmd struct {
	Name        string
	Version     string
	Description string
	Runner      Runner
	endpoint    chan string
	output      chan []byte
	validator   func()
	collector   func()
}

func Logger() zerolog.Logger {
	return logger
}

func NewCmd(name, version, description string, runner Runner) *Cmd {
	logger = logging.SetupLogger(name)
	return &Cmd{
		Name:        name,
		Version:     version,
		Description: description,
		Runner:      runner,
		endpoint:    make(chan string),
		output:      make(chan []byte),
	}
}

func Execute(cmd *Cmd) {
	cmd.check()

	rootCmd := cmd.newRootCmd()

	// endpoint flag
	rootCmd.PersistentFlags().StringP("endpoint", "e", "", "host:port notation of the endpoint to connect to")
	if err := rootCmd.MarkPersistentFlagRequired("endpoint"); err != nil {
		logger.Fatal().Err(err).Msg("Failed to mark flag required")
	}
	viper.BindPFlag("endpoint", rootCmd.PersistentFlags().Lookup("endpoint"))

	// endpoint output_folder flag
	rootCmd.PersistentFlags().StringP("output_folder", "o", "", "folder to store collected data")
	if err := rootCmd.MarkPersistentFlagRequired("output_folder"); err != nil {
		logger.Fatal().Err(err).Msg("Failed to mark flag required")
	}
	viper.BindPFlag("output_folder", rootCmd.PersistentFlags().Lookup("output_folder"))

	collectCmd := cmd.newCollectCmd()
	rootCmd.AddCommand(collectCmd)
	validateCmd := cmd.newValidateCmd()
	rootCmd.AddCommand(validateCmd)

	done = make(chan bool)

	go func() {
		defer close(done)
		if err := rootCmd.Execute(); err != nil {
			if ec, ok := err.(*ExitCodeError); ok {
				os.Exit(ec.Code)
			} else {
				os.Exit(1)
			}
		}
	}()
}

func Wait() {
	<-done
}

func (c *Cmd) newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s_%s", c.Name, c.Version),
		Short: c.Description,
		Long:  c.Description,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("a subcommand is required")
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return c.initConfig()
		},
		// SilenceErrors: true,
	}
}

func (c *Cmd) CallEndpoint(endpoint string) {
	c.endpoint <- endpoint
}

func (c *Cmd) ReadResult() []byte {
	return <-c.output
}

func (c *Cmd) RegisterValidator(f func()) {
	c.validator = f
}

func (c *Cmd) RegisterCollector(f func()) {
	c.collector = f
}

func (c *Cmd) check() {
	if c.Name == "" {
		logger.Fatal().Msg("Empty Name was provided")
	}
	if !semver.IsValid(fmt.Sprintf("v%s", c.Version)) {
		logger.Fatal().Msgf("Version does not follow Semantic Versioning format: %s", c.Version)
	}
}

func (c *Cmd) initConfig() error {
	logger.Info().Msgf("Initializing configuration for policy %s, version %s", c.Name, c.Version)
	logger.Info().Msg("Checking environment variables and flags")
	logger.Info().Msg("Checking if CLI_USER is set")
	err := viper.BindEnv("user", "CLI_USER")

	if err != nil {
		return fmt.Errorf("failed to bind CLI_USER: %w", err)
	}

	if viper.GetString("user") == "" {
		return fmt.Errorf("CLI_USER environment variable is required but not set")
	}
	logger.Info().Msgf("Using user %s", viper.GetString("user"))

	logger.Info().Msg("Checking if CLI_PASSWORD is set")
	err = viper.BindEnv("password", "CLI_PASSWORD")

	if err != nil {
		return fmt.Errorf("failed to bind CLI_PASSWORD: %w", err)
	}

	if viper.GetString("password") == "" {
		return fmt.Errorf("CLI_PASSWORD environment variable is required but not set")
	}
	logger.Info().Msg("Password is set")

	err = checkFolder(viper.GetString("output_folder"))
	if err != nil {
		return fmt.Errorf("failed to check output folder: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to check working directory: %w", err)
	}
	logger.Info().Msgf("Output folder is set to %s", filepath.Join(wd, viper.GetString("output_folder")))

	return nil
}

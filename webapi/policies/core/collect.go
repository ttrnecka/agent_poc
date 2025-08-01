package core

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
func (c Cmd) newCollectCmd() *cobra.Command {

	var cmd = &cobra.Command{
		Use:          "collect",
		Short:        "collect data",
		Long:         `collectdata`,
		RunE:         c.collect,
		SilenceUsage: true,
	}
	return cmd
}

func (c *Cmd) collect(cmd *cobra.Command, args []string) error {

	logger.Info().Msgf("Collecting data for policy %s, version %s", c.Name, c.Version)

	client, err := c.Runner.Connect()
	if err != nil {
		logger.Error().Err(err).Msg("Collection failed")
		return err
	}
	defer client.Close()

	if c.collector == nil {
		logger.Fatal().Msg("*collector* function is not defined")
	}

	go func() {
		c.collector()
		close(c.endpoint)
	}()

	exErr := &ExitCodeError{}
	// calls each command and saves each output to tagged file
	for cmd := range c.endpoint {
		out, err := c.Runner.Run(cmd)
		if err != nil {
			// save last error and code
			exErr = err
		}
		c.output <- out
	}
	if exErr.Code != 0 {
		logger.Error().Msgf("Collection failed, exit code: %d", exErr.Code)
		return exErr
	}
	logger.Info().Msg("Collection completed successfully")
	return nil
}

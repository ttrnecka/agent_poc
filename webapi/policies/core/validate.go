package core

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
func (c Cmd) newValidateCmd() *cobra.Command {

	var cmd = &cobra.Command{
		Use:          "validate",
		Short:        "validate access",
		Long:         `validate access`,
		RunE:         c.validate,
		SilenceUsage: true,
	}

	return cmd
}

func (c *Cmd) validate(cmd *cobra.Command, args []string) error {

	logger.Info().Str("policy", c.Name).Str("version", c.Version).Msg("Validating policy")

	client, err := c.Runner.Connect()
	if err != nil {
		logger.Error().Err(err).Msg("Validation failed")
		return err
	}
	defer client.Close()

	if c.validator == nil {
		logger.Fatal().Msg("*validator* function is not defined")
	}

	go func() {
		c.validator()
		close(c.endpoint)
	}()

	exErr := &ExitCodeError{}
	for cmd := range c.endpoint {
		out, err := c.Runner.Run(cmd)
		if err != nil {
			// save last error and code
			exErr = err
		}
		c.output <- out
	}
	if exErr.Code != 0 {
		logger.Error().Int("exit_code", exErr.Code).Msg("Validation failed")
		return exErr
	}
	logger.Info().Msg("Validation completed successfully")
	return nil
}

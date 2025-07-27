package core

import (
	"log"

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

	log.Printf("Validating probe %s version %s", c.Name, c.Version)

	client, err := c.Runner.Connect()
	if err != nil {
		log.Printf("Validation failed: %v", err)
		return err
	}
	defer client.Close()

	if c.validator == nil {
		log.Fatal("validator function not defined")
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
		log.Printf("Validation failed with exit code %d\n", exErr.Code)
		return exErr
	}
	log.Println("Validation completed successfully")
	return nil
}

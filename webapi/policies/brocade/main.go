package main

import (
	"github.com/rs/zerolog"
	"github.com/ttrnecka/agent_poc/webapi/policies/core"
)

var (
	NAME        = "brocade_cli"
	VERSION     = "1.0.1"
	DESCRIPTION = "Brocade CLI Collection Plugin"
)

var commandsV100 []string = []string{
	"version",
}

var commandsV101 []string = []string{
	"version",
	"switchshow",
}

var commandsV102 []string = []string{
	"version",
	"switchshow",
	"fabricshow",
	"licenseshow",
}

var commandsV103 []string = []string{
	"version",
	"switchshow",
	"fabricshow",
	"license --show",
}

var logger zerolog.Logger

func main() {
	cmd := core.NewCmd(NAME, VERSION, DESCRIPTION, &core.SshRunner{})
	logger = core.Logger()

	validate := func() {
		// all validation CMDs needs to have retrun code 0
		var validationCmds []string = []string{
			"version",
		}

		for _, endp := range validationCmds {
			// after every CallEndpoint we need to call ReadResult to prevent blocking
			cmd.CallEndpoint(endp)
			cmd.ReadResult()
			// here you can read result and do parsing if needed
			// no need to save anything, the core.Cmd already does the job
		}
	}

	collect := func() {
		// all validation CMDs needs to have retrun code 0
		var commands []string
		switch VERSION {
		case "1.0.0":
			commands = commandsV100
		case "1.0.1":
			commands = commandsV101
		case "1.0.2":
			commands = commandsV102
		case "1.0.3":
			commands = commandsV103
		default:
			logger.Fatal().Str("version", VERSION).Msg("Unknown version")
		}

		for _, endp := range commands {
			// after every CallEndpoint we need to call ReadResult to prevent blocking
			cmd.CallEndpoint(endp)
			cmd.ReadResult()
			// here you can read result and do parsing if needed
			// no need to save anything, the core.Cmd already does the job
		}
	}

	cmd.RegisterValidator(validate)
	cmd.RegisterCollector(collect)
	core.Execute(cmd)
	core.Wait()
}

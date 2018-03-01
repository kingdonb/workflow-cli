package parser

import (
	"github.com/teamhephy/workflow-cli/cmd"
	docopt "github.com/docopt/docopt-go"
)

// Version displays the client version
func Version(argv []string, cmdr cmd.Commander) error {
	usage := `
Displays the client version.

Usage: hephy version [options]

Options:
  -a --all
    list api and controller versions
`
	args, err := docopt.Parse(usage, argv, true, "", false, true)
	if err != nil {
		return err
	}

	return cmdr.Version(args["--all"].(bool))
}

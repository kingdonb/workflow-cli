package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	deis "github.com/teamhephy/controller-sdk-go"
	"github.com/teamhephy/workflow-cli/pkg/git"
	"github.com/teamhephy/workflow-cli/settings"
)

var defaultLimit = -1

func progress(wOut io.Writer) chan bool {
	frames := []string{"...", "o..", ".o.", "..o"}
	backspaces := strings.Repeat("\b", 3)
	tick := time.Tick(400 * time.Millisecond)
	quit := make(chan bool)
	go func() {
		for {
			for _, frame := range frames {
				fmt.Fprint(wOut, frame)
				select {
				case <-quit:
					fmt.Fprint(wOut, backspaces)
					close(quit)
					return
				case <-tick:
					fmt.Fprint(wOut, backspaces)
				}
			}
		}
	}()
	return quit
}

// load loads settings file and looks up the app name
func load(cf string, appID string) (*settings.Settings, string, error) {
	s, err := settings.Load(cf)

	if err != nil {
		return nil, "", err
	}

	if appID == "" {
		appID, err = git.DetectAppName(git.DefaultCmd, s.Client.ControllerURL.Host)

		if err != nil {
			return nil, "", err
		}
	}

	return s, appID, nil
}

func drinkOfChoice() string {
	drink := os.Getenv("DEIS_DRINK_OF_CHOICE")

	if drink == "" {
		drink = "coffee"
	}

	return drink
}

func limitCount(objs, total int) string {
	if objs == total {
		return "\n"
	}

	return fmt.Sprintf(" (%d of %d)\n", objs, total)
}

// checkAPICompatibility handles specific behavior for certain errors,
// such as printing an warning for the API mismatch error
func (d *DeisCmd) checkAPICompatibility(c *deis.Client, err error) error {
	if err == deis.ErrAPIMismatch {
		if !d.Warned {
			d.PrintErrf(`!    WARNING: Client and server API versions do not match. Please consider upgrading.
!    Client version: %s
!    Server version: %s
`, deis.APIVersion, c.ControllerAPIVersion)
			d.Warned = true
		}

		// API mismatch isn't fatal, so after warning continue on.
		return nil
	}

	return err
}

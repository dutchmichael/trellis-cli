package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mitchellh/cli"
	"trellis-cli/trellis"
)

type InitCommand struct {
	UI      cli.Ui
	Trellis *trellis.Trellis
}

func (c *InitCommand) Run(args []string) int {
	if err := c.Trellis.LoadProject(); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	switch len(args) {
	case 0:
	default:
		c.UI.Error(fmt.Sprintf("Error: too many arguments (expected 0, got %d)\n", len(args)))
		c.UI.Output(c.Help())
		return 1
	}

	c.UI.Info(fmt.Sprintf("Creating virtualenv in %s", c.Trellis.Virtualenv.Path))

	err := c.Trellis.Virtualenv.Create()
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating virtualenv: %s", err))
		return 1
	}

	c.UI.Info("✓ Virtualenv created\n")

	pip := execCommand("pip", "install", "-r", "requirements.txt")

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Installing pip dependencies (pip install -r requirements.txt) ..."
	s.FinalMSG = "✓ Dependencies installed\n"
	s.Start()

	err = pip.Run()

	if err != nil {
		s.Stop()
		c.UI.Error(fmt.Sprintf("Error installing pip requirements", err))
		return 1
	}

	s.Stop()

	return 0
}

func (c *InitCommand) Synopsis() string {
	return "Initializes an existing Trellis project"
}

func (c *InitCommand) Help() string {
	helpText := `
Usage: trellis init [options]

Initializes an existing Trellis project to be managed by trellis-cli.
The initialization process does two things:

1. creates a virtual environment via virtualenv to manage dependencies
2. installs dependencies via pip (specified by requirements.txt in your Trellis project)

This command is idempotent meaning it can be run multiple times without side-effects.

  $ trellis init

Options:
  -h, --help  show this help
`

	return strings.TrimSpace(helpText)
}

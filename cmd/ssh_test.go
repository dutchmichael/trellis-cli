package cmd

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"trellis-cli/trellis"
)

func TestSshRunValidations(t *testing.T) {
	defer trellis.LoadFixtureProject(t)()
	ui := cli.NewMockUi()

	cases := []struct {
		name            string
		projectDetected bool
		args            []string
		out             string
		code            int
	}{
		{
			"no_project",
			false,
			nil,
			"No Trellis project detected",
			1,
		},
		{
			"no_args",
			true,
			nil,
			"Usage: trellis",
			1,
		},
		{
			"invalid_env",
			true,
			[]string{"foo"},
			"Error: foo is not a valid environment",
			1,
		},
		{
			"invalid_env_with_site",
			true,
			[]string{"foo", "example.com"},
			"Error: foo is not a valid environment",
			1,
		},
		{
			"invalid_site",
			true,
			[]string{"development", "nosite"},
			"Error: nosite is not a valid site",
			1,
		},
		{
			"too_many_args",
			true,
			[]string{"development", "site", "foo"},
			"Error: too many arguments",
			1,
		},
	}

	for _, tc := range cases {
		mockProject := &MockProject{tc.projectDetected}
		trellis := trellis.NewTrellis(mockProject)
		sshCommand := &SshCommand{ui, trellis}

		code := sshCommand.Run(tc.args)

		if code != tc.code {
			t.Errorf("expected code %d to be %d", code, tc.code)
		}

		combined := ui.OutputWriter.String() + ui.ErrorWriter.String()

		if !strings.Contains(combined, tc.out) {
			t.Errorf("expected output %q to contain %q", combined, tc.out)
		}
	}
}
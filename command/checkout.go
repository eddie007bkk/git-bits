package command

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"github.com/nerdalize/git-bits/bits"
)

type Checkout struct {
	ui cli.Ui
}

func NewCheckout() (cmd cli.Command, err error) {
	return &Checkout{
		ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stderr,
			ErrorWriter: os.Stderr,
		},
	}, nil
}

// Help returns long-form help text that includes the command-line
// usage, a brief few sentences explaining the function of the command,
// and the complete list of flags the command accepts.
func (cmd *Checkout) Help() string {
	return fmt.Sprintf(`
  %s
`, cmd.Synopsis())
}

// Synopsis returns a one-line, short synopsis of the command.
// This should be less than 50 characters ideally.
func (cmd *Checkout) Synopsis() string {
	return "fetch chunks and combine them into the orignal file"
}

// Run runs the actual command with the given CLI instance and
// command-line arguments. It returns the exit status when it is
// finished.
func (cmd *Checkout) Run(args []string) int {
	wd, err := os.Getwd()
	if err != nil {
		cmd.ui.Error(fmt.Sprintf("Failed to get working directory: %v", err))
		return 1
	}

	repo, err := bits.NewRepository(wd, os.Stderr)
	if err != nil {
		cmd.ui.Error(fmt.Sprintf("failed to setup repository: %v", err))
		return 2
	}

	lstore, err := repo.LocalStore()
	if err != nil {
		cmd.ui.Error(fmt.Sprintf("failed to open local store: %v", err))
		return 3
	}

	defer lstore.Close()
	err = repo.Checkout(lstore, os.Stdin, os.Stdout)
	if err != nil {
		cmd.ui.Error(fmt.Sprintf("failed to combine: %v", err))
		return 4
	}

	return 0
}
// Package ui is the dot CLI's interactive layer: a small, testable Prompter
// abstraction backed by charmbracelet/huh, run against /dev/tty so prompts reach
// the real terminal even when the caller has captured stdout.
package ui

import (
	"errors"
	"os"

	"github.com/charmbracelet/huh"
)

// ErrNoTTY is returned by the real Prompter when there is no controlling
// terminal, so callers can fall back (e.g. skip a confirmation with a warning).
var ErrNoTTY = errors.New("no controlling terminal")

// Prompter asks the user interactive questions.
type Prompter interface {
	// Confirm asks a yes/no question, returning def if the user aborts (Esc/Ctrl-C).
	// It returns ErrNoTTY when there is no terminal to prompt on.
	Confirm(title string, def bool) (bool, error)
	// Select asks the user to pick one option, returning "" if they abort.
	Select(title string, options []string) (string, error)
	// MultiSelect asks the user to pick any number of options, returning nil if
	// they abort or pick nothing.
	MultiSelect(title string, options []string) ([]string, error)
}

// huhPrompter is the real Prompter: charmbracelet/huh on /dev/tty.
type huhPrompter struct{}

// New returns a Prompter backed by huh on the controlling terminal.
func New() Prompter { return huhPrompter{} }

func (huhPrompter) Confirm(title string, def bool) (bool, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return def, ErrNoTTY
	}
	defer tty.Close()

	value := def
	form := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().Title(title).Affirmative("Yes").Negative("No").Value(&value),
	)).WithInput(tty).WithOutput(tty)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return def, nil
		}
		return def, err
	}
	return value, nil
}

func (huhPrompter) Select(title string, options []string) (string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return "", ErrNoTTY
	}
	defer tty.Close()

	var choice string
	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title(title).Options(huh.NewOptions(options...)...).Value(&choice).Filtering(true),
	)).WithInput(tty).WithOutput(tty)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return "", nil
		}
		return "", err
	}
	return choice, nil
}

func (huhPrompter) MultiSelect(title string, options []string) ([]string, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil, ErrNoTTY
	}
	defer tty.Close()

	var chosen []string
	form := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[string]().Title(title).Options(huh.NewOptions(options...)...).Value(&chosen).Filterable(true),
	)).WithInput(tty).WithOutput(tty)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return nil, nil
		}
		return nil, err
	}
	return chosen, nil
}

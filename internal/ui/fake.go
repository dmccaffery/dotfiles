package ui

// Fake is a scripted Prompter for tests. Each queue is consumed in order; when a
// queue is exhausted the zero value is returned (no prompt).
type Fake struct {
	Replies         []bool     // Confirm answers
	Selections      []string   // Select answers
	MultiSelections [][]string // MultiSelect answers
	// NoTTY makes every method return ErrNoTTY (a non-interactive context).
	NoTTY bool
	// Asked records each prompt title, in order.
	Asked []string
}

// Confirm implements Prompter.
func (f *Fake) Confirm(title string, def bool) (bool, error) {
	f.Asked = append(f.Asked, title)
	if f.NoTTY {
		return def, ErrNoTTY
	}
	if len(f.Replies) == 0 {
		return def, nil
	}
	r := f.Replies[0]
	f.Replies = f.Replies[1:]
	return r, nil
}

// Select implements Prompter.
func (f *Fake) Select(title string, _ []string) (string, error) {
	f.Asked = append(f.Asked, title)
	if f.NoTTY {
		return "", ErrNoTTY
	}
	if len(f.Selections) == 0 {
		return "", nil
	}
	s := f.Selections[0]
	f.Selections = f.Selections[1:]
	return s, nil
}

// MultiSelect implements Prompter.
func (f *Fake) MultiSelect(title string, _ []string) ([]string, error) {
	f.Asked = append(f.Asked, title)
	if f.NoTTY {
		return nil, ErrNoTTY
	}
	if len(f.MultiSelections) == 0 {
		return nil, nil
	}
	s := f.MultiSelections[0]
	f.MultiSelections = f.MultiSelections[1:]
	return s, nil
}

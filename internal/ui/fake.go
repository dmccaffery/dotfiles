package ui

// Fake is a scripted Prompter for tests.
type Fake struct {
	// Replies are consumed in order by Confirm; once exhausted, def is returned.
	Replies []bool
	// NoTTY makes Confirm return ErrNoTTY (simulating a non-interactive context).
	NoTTY bool
	// Asked records each Confirm title, in order.
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

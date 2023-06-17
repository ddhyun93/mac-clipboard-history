package tui

type Item string

func (i Item) FilterValue() string {
	return string(i)
}

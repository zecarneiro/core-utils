package cobralib

type FlagsOptions[T any] struct {
	Name       string
	Shorthand  string
	Usage      string
	IsRequired bool
	Default    T
}

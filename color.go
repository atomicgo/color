package color

type Color interface {
	Sequence(background bool) string
}

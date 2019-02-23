package model

type Page struct {
	Number int
	Size   int
}

func (p Page) StartIndex() int {
	return (p.Number - 1) * p.Size
}

func (p Page) EndPosition() int {
	return p.Number * p.Size
}

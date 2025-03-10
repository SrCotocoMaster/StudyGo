package main

type car struct {
	Make   string
	Model  string
	Height int
	Width  int
}

func (c car) getHorsePower() int {
	panic("implement me")
}

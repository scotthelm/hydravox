package main

type ApiRoot struct {
	Meta `json:"meta"`
}
type Meta struct {
	Name      string `json:"name"`
	Licensing string `json:"licensing"`
}

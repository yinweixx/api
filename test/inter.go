package test

import (
	"fmt"
)

type YW interface {
	Get(i string) error

	Put()
}

type mime struct {
	i string
}

func (m *mime) Get(i string) error {
	fmt.Println(i)
	return nil
}

func (m *mime) Put() {
}

// func (m mime) Get(i string) error {
// 	fmt.Println(i)
// 	return nil
// }

// func (m mime) Put() {}

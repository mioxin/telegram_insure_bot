package services

import "fmt"

type Insurance struct {
	Name  string
	Price float64
}

func NewInsurence(name string, price float64) *Insurance {
	return &Insurance{name, price}
}

func (ins *Insurance) Calculate(data any) (string, error) {
	var sum string = fmt.Sprintf("%.2f тенге", ins.Price)
	return sum, nil
}

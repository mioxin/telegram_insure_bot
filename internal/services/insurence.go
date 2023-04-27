package services

import "fmt"

type Insurance struct {
	Name                                string
	Total_work, Vid, Workers1, Workers2 int
	Gfot1, Gfot2                        float64
}

func NewInsurence(name string, price float64) *Insurance {
	return &Insurance{Name: name}
}

func (ins *Insurance) Calculate() (string, error) {
	var sum string = fmt.Sprintf("%.2f тенге", 10.0)
	return sum, nil
}

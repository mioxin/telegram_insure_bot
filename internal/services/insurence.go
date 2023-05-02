package services

import (
	"fmt"
	"log"
	"time"
)

type Insurance struct {
	Name                                string
	Total_work, Vid, Workers1, Workers2 int
	Gfot1, Gfot2                        float64
}

func NewInsurence(name string, price float64) *Insurance {
	return &Insurance{Name: name}
}

func (ins *Insurance) Calculate() (string, error) {
	sum := 10.0
	bonus := 1.5

	//TODO calculate sum
	log.Println("Calculate 5 sec ...")
	time.Sleep(5 * time.Second)
	var str string = fmt.Sprintf("Сумма страховки: *%.2f тенге*\nВаша скидка: *%.2f тенге*", sum, bonus)
	return str, nil
}

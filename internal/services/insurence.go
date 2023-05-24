package services

import (
	"fmt"
	"log"
	"time"
)

type Insurance struct {
	Name, VidDescr, BinIin, Vid string
	Total_work                  int
	Gfot                        float64
	EventInLast5Year            bool
}

func NewInsurence(name string) *Insurance {
	return &Insurance{Name: name}
}

func (ins *Insurance) Calculate() (string, error) {
	sum := 70000.0
	bonus := 23000.0

	//TODO calculate sum
	log.Println("Calculate 5 sec ...")
	time.Sleep(1 * time.Second)
	var str string = fmt.Sprintf("Сумма страховки: *%.2f тенге*\nВаша скидка: *%.2f тенге*", sum, bonus)
	return str, nil
}

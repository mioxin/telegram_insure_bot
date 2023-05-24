package services

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type fileTypeOfBusns map[string]string

func NewFileTypeOfBusns(file string) (fileTypeOfBusns, error) {
	tob := make(map[string]string)
	strInp, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err != nil {
		log.Println("error NewFileTypeOdBusns: Cant open file vid.txt.", err)
		return nil, err
	}
	in := bufio.NewReader(strings.NewReader(string(strInp)))
	for {
		ln, err := in.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("error NewFileTypeOdBusns: error read vid.txt.", err)
			return nil, err
		}
		arr := strings.Split(strings.TrimSpace(ln), "-")
		vid := strings.TrimSpace(arr[0])
		if len(arr) < 2 {
			log.Printf("error NewFileTypeOdBusns: skip line %v.\n %v", ln, err)
		} else {
			tob[vid] = arr[1]
		}
	}
	return fileTypeOfBusns(tob), nil
}

func (f fileTypeOfBusns) Get(vid string) (string, error) {
	if descr, ok := f[vid]; ok {
		return descr, nil
	} else {
		return "", fmt.Errorf("not fount type og business %v", vid)
	}
}

// func str2int(str string) (int, error) {
// 	res, err := strconv.Atoi(strings.TrimSpace(str))
// 	return res, err
// }

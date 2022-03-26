package main

import (
	"fmt"
	"os"
)

var Number uint

func init() {
	Number = 0
}

func main() {
	mydir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mydir)
	Number += 1
	fmt.Println(Number)
}

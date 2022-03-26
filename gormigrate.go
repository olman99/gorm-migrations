package main

import (
	"fmt"
	"path"
)

func main() {
	fmt.Println(path.Base("."))
	fmt.Println(path.Dir("."))
}

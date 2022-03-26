package main

import (
	"log"
	"path"
	"runtime"
)

func main() {
	if _, file, _, ok := runtime.Caller(0); ok {
		__dirname := path.Dir(file)
		log.Println("__dirname:", __dirname)
	}
}

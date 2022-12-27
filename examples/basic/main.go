package main

import (
	"fmt"
	"github.com/bxcodec/gotcha"
	"log"
)

func main() {
	err := gotcha.Set("name", "John Snow")
	if err != nil {
		log.Fatal(err)
	}
	val, err := gotcha.Get("name")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)

	//Output: John Snow
}

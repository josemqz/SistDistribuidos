package main

import (
	"log"
	"math/rand"
)

func main(){

	for i := 0; i < 10; i++{
		log.Println(rand.Intn(3))
	}
}

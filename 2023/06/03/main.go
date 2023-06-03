package main

import (
	"fmt"
	"time"
)

func main() {
	water := putWater()
	tea := putTea()
	fmt.Println("waiting for tea and water")

	fmt.Println("do something else")
	time.Sleep(10 * time.Second)
	w := <-water
	t := <-tea
	fmt.Println("tea and water are ready", w, t)

}

func putWater() <-chan string {
	water := make(chan string)
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("water is ready")
		water <- "water"
	}()
	return water
}

func putTea() <-chan string {
	tea := make(chan string)
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("tea is ready")
		tea <- "tea"
	}()
	return tea
}

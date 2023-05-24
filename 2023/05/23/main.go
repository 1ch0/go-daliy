package main

import (
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"
)

func main() {
	demo01()
	demo02()
}

func demo01() {
	name := pie.Of([]string{"Bob", "Sally", "John", "Jane"}).
		FilterNot(func(name string) bool {
			return strings.HasPrefix(name, "J")
		}).Map(strings.ToUpper).Last()

	fmt.Println(name)
}

func demo02() {
	slice01 := []int{1, 2, 3}
	slice02 := []string{"a", "b", "c"}

	fmt.Printf("slice01 are sorted: %v\n", pie.AreSorted(slice01))
	fmt.Printf("slice02 are sorted: %v\n", pie.AreSorted(slice02))
}

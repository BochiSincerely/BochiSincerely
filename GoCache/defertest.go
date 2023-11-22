package main

import "fmt"

func main() {
	value := 0

	defer func() {
		fmt.Println("defer:", value)
	}()

	value++
	fmt.Println("after increment:", value)

}

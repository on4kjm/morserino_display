package main

import (
	"fmt"
	"time"
)

func main() {
	const name, age = "Kim", 22
	fmt.Print(name, " is ", age, " years old.\n")
	fmt.Print("1 ")
	time.Sleep(2 * time.Second)
	fmt.Print("2 ")
	time.Sleep(2 * time.Second)
	fmt.Print("3 ")
	time.Sleep(2 * time.Second)
	fmt.Print("4 ")
	time.Sleep(2 * time.Second)
	fmt.Print("\n")

	fmt.Print("1 ")
	time.Sleep(2 * time.Second)
	fmt.Print("\x1b[2D")
	fmt.Print("2 ")
	time.Sleep(2 * time.Second)
	fmt.Print("\x1b[2D")
	fmt.Print("3 ")
	time.Sleep(2 * time.Second)
	fmt.Print("\x1b[2D")
	fmt.Print("4 ")
	time.Sleep(2 * time.Second)
	fmt.Print("\n")
	

	// It is conventional not to worry about any
	// error returned by Print.

}
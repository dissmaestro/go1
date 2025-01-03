package main

import "fmt"

// Hi returns a friendly greeting
func Hi(name string) string {
	return fmt.Sprintf("Hi, %s", name)
}

func main() {
	fmt.Println(Hi("dosya"))
	fmt.Println("Kakogo huya")

}

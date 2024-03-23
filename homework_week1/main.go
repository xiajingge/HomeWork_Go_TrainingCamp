package main

import "fmt"

func main() {
	s1 := make([]int, 0, 4)
	s1 = append(s1, 1, 2, 3)
	println("Before deletion: ")
	for _, v := range s1 {
		println(v)
	}
	s1, i, err := Delete(s1, 33)
	if err != nil {
		println("err", err)
		return
	} else {
		println("After deletion: ")
		for _, v := range s1 {
			println(v)
		}
		fmt.Printf("Deleted element: %d", i)
	}
}

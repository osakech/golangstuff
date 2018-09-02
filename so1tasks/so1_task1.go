package main

import (
	"fmt"
)

func task1solution(A []int) int {
	// write your code in Go 1.4
	first := A[0]
	counter := 1
	if first == -1 {
		return counter
	}
	llist(A, first, &counter)
	return counter
}

func llist(A []int, idx int, c *int) {
	*c++
	if A[idx] == -1 {
		return
	}
	llist(A, A[idx], c)
	return
}

func runtask1() {
	//         0, 1, 2, 3, 4,  5, 6, 7, 8
	a := []int{2, 4, 1, 5, 8, -1, 5, 6, 5}
	counter := task1solution(a)
	fmt.Println(counter)
}

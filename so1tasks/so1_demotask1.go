package main

// This is a demo task.

// Write a function:

// func Solution(A []int) int

// that, given an array A of N integers, returns the smallest positive integer (greater than 0) that does not occur in A.

// For example, given A = [1, 3, 6, 4, 1, 2], the function should return 5.

// Given A = [1, 2, 3], the function should return 4.

// Given A = [−1, −3], the function should return 1.

// Assume that:

// N is an integer within the range [1..100,000];
// each element of array A is an integer within the range [−1,000,000..1,000,000].
// Complexity:

// expected worst-case time complexity is O(N);
// expected worst-case space complexity is O(N) (not counting the storage required for input arguments).

// you can also use imports, for example:
// import "fmt"
// import "os"

// you can write to stdout for debugging purposes, e.g.
// fmt.Println("this is a debug message")

import (
	"fmt"
	"sort"
)

func rundemotask() {
	ints := []int{1, 3, 6, 4, 1, 2}
	s := demo(ints)
	fmt.Println(s)
}

func demo(A []int) int {
	var b int
	sort.Ints(A)

	for i := range A {
		nextInt := A[i] + 1
		if A[i+1] > nextInt {
			return nextInt
		}
		b = A[i+1]
	}

	if b <= 0 {
		return 1
	}

	return b

}

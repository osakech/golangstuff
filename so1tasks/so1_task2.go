package main

// failed this one

// type Tree struct {
// 	X int
// 	L *Tree
// 	R *Tree
// }

// func main() {

// 	// res := plus(1, 2)

// 	a := &Tree{
// 		X: 5,
// 		L: &Tree{
// 			X: 3,
// 			L: &Tree{
// 				X: 20,
// 				L: nil,
// 				R: nil,
// 			},
// 			R: &Tree{
// 				X: 21,
// 				L: nil,
// 				R: nil,
// 			},
// 		},
// 		R: &Tree{
// 			X: 10,
// 			L: &Tree{
// 				X: 1,
// 				L: nil,
// 				R: nil,
// 			},
// 			R: nil,
// 		},
// 	}
// 	res := solution(a)
// 	fmt.Println(res)

// }

// // func plus(a int, b int) int {
// // 	return a + b
// // }
// //(5, (3, (20, None, None), (21, None, None)), (10, (1, None, None), None))
// //(5, (3, (20, None, None), (21, None, None)), (10, (1, None, None), None))

// func solution(T *Tree) int {
// 	// write your code in Go 1.4
// 	c := make(chan int)
// 	increment := 0
// 	rec(T, c)
// 	increment += <-c
// 	return increment
// }

// func rec(T *Tree, c chan int) {
// 	if T.L == nil && T.R == nil {
// 		c <- 1
// 	}
// 	if T.L != nil {
// 		go rec(T.L, c)
// 	}
// 	if T.R != nil {
// 		go rec(T.R, c)
// 	}

// 	return
// }

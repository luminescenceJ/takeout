package main

import (
	"fmt"
	"sort"
)

func main() {
	// a := 0
	// b := 0
	// for {
	//     n, _ := fmt.Scan(&a, &b)
	//     if n == 0 {
	//         break
	//     } else {
	//         fmt.Printf("%d\n", a + b)
	//     }
	// }
	var n int
	_, _ = fmt.Scan(&n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		_, _ = fmt.Scan(&arr[i])
	}
	sort.Ints(arr)
	for i := 0; i < n; {
		if i+1 >= len(arr) {
			break
		}
		if arr[i] == arr[i+1] {
			arr = append(arr[:i], arr[i+1:]...)
		} else {
			i++
		}

	}

	for _, n := range arr {
		fmt.Println(n)
	}

}

package main

import "fmt"

func wardrobeFinishing(m int, n int, cnt int) int {
	var res int
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			// 计算 i 和 j 的各自数位之和
			sum := digit(i) + digit(j)
			// 如果和小于等于 cnt，则计数
			if sum <= cnt {
				res++
			}
		}
	}
	return res
}

func digit(x int) int {
	result := 0
	// 计算各位数字之和
	for x > 0 {
		result += x % 10
		x = x / 10
	}
	return result
}

func main() {
	m := 16
	n := 8
	cnt := 4
	fmt.Println(wardrobeFinishing(m, n, cnt)) // 预期输出 15
}

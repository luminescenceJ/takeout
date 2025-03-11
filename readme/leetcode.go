package main

import (
	"bufio"
	"fmt"
	"os"
)

type loc struct {
	x int
	y int
}

func check(matrix [][]byte, dst loc, key map[byte]bool) bool {
	if dst.x >= len(matrix) || dst.x < 0 || dst.y < 0 || dst.y >= len(matrix[0]) || matrix[dst.x][dst.y] == 0 {
		return false
	} else if matrix[dst.x][dst.y] <= 'Z' && matrix[dst.x][dst.y] >= 'A' {
		if key[byte(matrix[dst.x][dst.y]-32)] {
			return true
		}
		return false
	}
	return true
}

func bfs(start, end loc, matrix [][]byte, key map[byte]bool, distance [][]int, step int) {
	// distance 维护到当前为止的最少步数，如果到这里的步数大于distance[i,j]表示绕圈，直接return
	if distance[start.x][start.y] == -1 {
		distance[start.x][start.y] = step
	} else if distance[start.x][start.y] <= step {
		return
	}

	if start == end {
		return
	}

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			next := loc{
				x: start.x + i,
				y: start.y + j,
			}
			if i == 0 && j == 0 || !check(matrix, next, key) {
				continue
			}

			if matrix[start.x][start.y] < 'z' && matrix[start.x][start.y] > 'a' {
				key[matrix[start.x][start.y]] = true
			}
			bfs(next, end, matrix, key, distance, step+1)
			if matrix[start.x][start.y] < 'z' && matrix[start.x][start.y] > 'a' {
				key[matrix[start.x][start.y]] = false
			}

		}
	}

}

func main() {
	n, m := 0, 0 // 行数和列数
	fmt.Scan(&n, &m)
	var start, end loc
	matrix := make([][]byte, n)

	inputs := bufio.NewScanner(os.Stdin)
	for i := 0; i < n; i++ {
		matrix[i] = make([]byte, m)
		inputs.Scan()
		s := inputs.Text()
		for j := range s {
			matrix[i][j] = s[j]
			if matrix[i][j] == '2' {
				start.x = j
				start.y = i
			} else if matrix[i][j] == '3' {
				end.x = j
				end.y = i
			}
		}
	}

	key := map[byte]bool{}
	distance := make([][]int, n)
	for i := range distance {
		distance[i] = make([]int, m)
		for j := range distance[i] {
			distance[i][j] = -1
		}
	}

	bfs(start, end, matrix, key, distance, 0)

	fmt.Println(distance[end.y][end.x])

}

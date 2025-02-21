package sort_algorithm

import "math/rand"

// 冒泡排序算法
func BubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

// 快速排序算法入口
func QuickSort(arr []int) {
	quickSort(arr, 0, len(arr)-1)
}

// 快排算法
func quickSort(arr []int, low, high int) {
	if low < high {
		pivot := partition(arr, low, high)
		quickSort(arr, low, pivot-1)
		quickSort(arr, pivot+1, high)
	}
}

// 快排拆分逻辑
func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// 生成指定长度的随机数字切片
func makeRandomNumberSlice(n int) []int {
	numbers := make([]int, n)
	for i := range numbers {
		numbers[i] = rand.Intn(n)
	}
	return numbers
}

const LENGTH = 10_000

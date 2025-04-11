package main

import "fmt"

func search(nums []int, target int) int {
	// 4,5,6,0,1,2 t=0
	// 二分，
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) / 2
		if nums[mid] == target {
			return mid
		}
		if nums[mid] < target {
			if nums[mid] < nums[right] && nums[right] >= target {
				left = mid + 1
			} else {
				right = mid - 1
			}
		} else {
			if nums[mid] > nums[right] && nums[left] < target {
				// left - mid 有序
				right = mid - 1
			} else {
				left = mid + 1
			}
		}
	}
	return -1
}

func main() {
	fmt.Println(search([]int{1, 3}, 3))
}

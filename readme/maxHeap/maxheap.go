package maxHeap

import "fmt"

type heap struct {
	arr []int
} // 大顶堆

func (h *heap) Push(x int) {
	h.arr = append(h.arr, x)
	h.heapifyUp(len(h.arr) - 1)
}

func (h *heap) Pop() int {
	if len(h.arr) == 0 {
		return -1
	}
	top := h.arr[0]
	h.arr[0], h.arr[len(h.arr)-1] = h.arr[len(h.arr)-1], h.arr[0]
	h.arr = h.arr[:len(h.arr)-1]
	h.heapifyDown(0)
	return top
}

func (h *heap) heapifyDown(index int) {
	n := len(h.arr)
	// 从上往下交换
	for {
		left := index*2 + 1
		right := index*2 + 2
		biggest := index
		if left < n && h.arr[left] > h.arr[biggest] {
			biggest = left
		}
		if right < n && h.arr[right] > h.arr[biggest] {
			biggest = right
		}
		if biggest == index {
			break
		}
		h.arr[biggest], h.arr[index] = h.arr[index], h.arr[biggest]
		index = biggest
	}
}

func (h *heap) heapifyUp(index int) {
	parent := (index - 1) / 2
	for index > 0 && h.arr[parent] < h.arr[index] {
		h.arr[index], h.arr[parent] = h.arr[parent], h.arr[index]
		index = parent
		parent = (index - 1) / 2
	}
}

func initHeap(arr []int) *heap {
	h := &heap{arr: arr}
	// 从第一个非叶子节点开始向下heapify
	for i := len(h.arr)/2 - 1; i >= 0; i-- {
		h.heapifyDown(i)
	}
	return h
}

func main() {
	h := initHeap([]int{1, 2, 3, 8, 9, 7, 1})
	fmt.Println(h.arr)
	for len(h.arr) > 0 {
		fmt.Println(h.Pop())
	}
}

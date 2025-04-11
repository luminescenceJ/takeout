package sort_algorithm

func bubbleSort(arr []int) []int {
	flag := false
	for i := 0; i < len(arr)-1; i++ {
		for j := 0; j < len(arr)-1-i; j++ {
			if arr[j] > arr[j+1] {
				flag = true
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
		if flag == false {
			break
		}
	}
	return arr
}

package sort_algorithm

import "testing"

// 基准测试运行命令: go test -bench=.   这里的.表示运行当前所有的基准测试, 也可以指定函数名
// benchmark基准测试用例

func BenchmarkBubbleSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer() // 停止计时
		numbers := makeRandomNumberSlice(LENGTH)

		b.StartTimer() // 开始计时
		BubbleSort(numbers)
	}
}

func BenchmarkQuickSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer() // 停止计时
		numbers := makeRandomNumberSlice(LENGTH)

		b.StartTimer() // 开始计时
		QuickSort(numbers)
	}
}

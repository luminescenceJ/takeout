package main

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func generateTestData(length int) []int {
	arr := make([]int, length)
	for i := 0; i < length; i++ {
		arr[i] = rand.Intn(1000)
	}
	return arr
}

func TestMergeSort(t *testing.T) {
	nums := generateTestData(1000)
	testNums := MergeSort(nums)
	sort.Ints(nums)
	for i := 0; i < len(nums); i++ {
		if !reflect.DeepEqual(nums, testNums) {
			t.Errorf("merge sort failed")
		}
	}
	t.Log("merge sort success")
}

func BenchmarkMergeSort(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testNums := generateTestData(1000)
		MergeSort(testNums)
	}
}

func BenchmarkSort(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testNums := generateTestData(1000)
		sort.Ints(testNums)
	}
}

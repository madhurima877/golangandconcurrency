package main

func Add(a, b int) int {
	return a + b
}
func Multiply(a, b int) int {
	return a * b
}
func SlowFunction(n int) int {
	result := 0
	for i := 0; i < n; i++ {
		result += i * i
	}
	return result
}
func AllocateMemory() []int {
	data := make([]int, 1000000)
	for i := range data {
		data[i] = i
	}
	return data
}

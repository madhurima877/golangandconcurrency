package main

import "log"

func main() {
	arr := []int{1, 2}
	appends(&arr, 3)
	log.Println(arr)
}
func appends(arr *[]int, i int) {
	*arr = append(*arr, i)
}

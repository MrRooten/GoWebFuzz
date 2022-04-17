package main

import (
	"fmt"
	"time"
)

func min[T int | float64](x, y T) T {

	if x < y {
		return x
	}
	return y
}

// 调用泛型函数
func test() {
	go func() {
		for ; ; {
			time.Sleep(1 * time.Second)
			fmt.Println("hello")
		}
	}()
	return
}
func main() {
	test()
	for ; ; {

	}
}

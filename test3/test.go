package main

import "gowebfuzz/library/utils"

func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}
func main() {
	config := utils.Config{}
	config.Initialize("./config.inc")
}

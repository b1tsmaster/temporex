package main

import "fmt"

func main() {
	v := make([]int, 5, 10)
	for i := range v {
		v[i] = i * i
	}
	v = append(v, 25)
	fmt.Printf("%v\n", v)
}

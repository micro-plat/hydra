package main

import (
	"fmt"
)

//	_ "github.com/micro-plat/hydra/hydra"

func main() {
	idx := []int{1, 2, 3, 4, 5, 6}
	i := 1
	nsubs := make([]int, 0, len(idx))
	nsubs = append(nsubs, idx[0:i]...)
	nsubs = append(nsubs, idx[i+1:]...)
	fmt.Println(nsubs)
}

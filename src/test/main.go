package main

import (
	"fmt"
	"net/rpc"
)

type T struct {
	t int
}

func (t *T) Do(i int) int {
	return t.t + i
}

type Do interface {
	Do(i int) int
}

func main() {
	rpc.Server()
	i := 33 - 24
	fmt.Println((i * (i - 1) * (i - 2) * (i - 3) * (i - 4) * (i - 5)) / (6 * 5 * 4 * 3 * 2) * 16)

}

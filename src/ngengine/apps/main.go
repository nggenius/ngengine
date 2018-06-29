package main

import (
	"fmt"
	"time"
)

var (
	uuid          int
	magic_time, _ = time.Parse("2006-01-02 15:04:05", "2018-01-01 00:00:00")
)

func GenerateGUID(id int) uint64 {
	uuid++
	dur := time.Now().Sub(magic_time).Seconds()
	ms := int(dur*10) - int(dur)*10
	if ms == 0 {
		ms = 1
	}
	return (uint64(id)&0xFFFF)<<48 |
		(uint64(dur)&0xFFFFFFFF)<<16 |
		(uint64(uuid%0xFFF)&0xFFF)<<4 |
		uint64(0xF/ms)&0xF
}

func main() {
	for i := 0; i < 4000; i++ {
		fmt.Printf("%X\n", GenerateGUID(1))
	}
}

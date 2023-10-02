package main

import (
	"fmt"
	"log"
)

func fn(m *map[int]int) {
	//mNew := make(map[int]int)
	*m = make(map[int]int)
	fmt.Println("m == nil in fn?:", m == nil)
}

func main() {
	//m := map[int]int{1: 1, 2: 2, 3: 3}
	//log.Println(1, len(m))
	//log.Println(1, m)
	//
	//delete(m, 1)
	//
	//log.Println(2, len(m))
	//log.Println(2, m)
	var s0 []string
	log.Println(len(s0))

	s1 := make([]int, 0, 4)
	log.Println("s1 v:", s1)
	log.Printf("s1 p: %p", s1)

	s1 = append(s1, 2)
	log.Println("s1 v:", s1)
	log.Printf("s1 p: %p", s1)

	//fn(&m)
	//fmt.Println("m == nil in main?:", m == nil)
}

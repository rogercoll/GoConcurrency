package main

import (
	"fmt"
	"sync"
)

func main() {
	aux := 0
	myPool := &sync.Pool{
		New: func() interface{} {
			aux += 1
			fmt.Println("Creating new instance")
			return aux
		},
	}

	instance := myPool.Get()
	instance2 := myPool.Get()
	myPool.Put(instance)
	myPool.Put(instance2)
	a1 := myPool.Get()
	a2 := myPool.Get()
	a3 := myPool.Get()
	fmt.Println(a1)
	fmt.Println(a2)
	fmt.Println(a3)
}

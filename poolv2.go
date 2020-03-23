package main

import (
	"fmt"
	"sync"
)

func main() {
	//Go has a garbage collector taht the instatntied objecteds will be automatically cleaned up
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	//Seed the pool with 4KB
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()
			mem := calcPool.Get().(*[]byte) //get the ADDRESS
			defer calcPool.Put(mem)

			//DO SOMETHING INTO THIS MEMORY SPACE ;)
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created\n", numCalcsCreated)

	//If we hadn't used a pool we would have allocated a gigabyte of memory, but we have allocated ONLY 4KB
	//POOL is a cache of preallocated objects for operations that must run as quickly as possible
}

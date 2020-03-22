package main

import (
	"fmt"
	"sync"
)

//Mutex => mutual exclusion

func main() {
	var count int
	var lock sync.Mutex

	increment := func() {
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("Incrementing: %d\n", count)
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("Decrementing: %d\n", count)
	}

	var arithmeticWG sync.WaitGroup
	for i := 0; i <= 5; i++ {
		arithmeticWG.Add(1)
		go func() {
			defer arithmeticWG.Done()
			increment()
		}()
	}

	for i := 0; i <= 5; i++ {
		arithmeticWG.Add(1)
		go func() {
			defer arithmeticWG.Done()
			decrement()
		}()
	}

	arithmeticWG.Wait()
	fmt.Println("Arithmetic operations complete")
}

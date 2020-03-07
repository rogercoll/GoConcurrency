package main

import (
	"sync"
	"testing"
)

//go test -bench=. -cpu=1 messageTimeThreads_test.go
//The benchmark function must run the target code b.N times. During benchmark execution, b.N is adjusted until the benchmark function lasts long enough to be timed reliably
// ns/op => mitjana dels segons per loop
//it takes 130 ns to share a struct(message) to another thread

func BenchmarkContextSwitch(b *testing.B) {
	var wg sync.WaitGroup
	begin := make(chan struct{})
	c := make(chan struct{})

	var message struct{}
	sender := func() {
		defer wg.Done()
		<-begin //wait until begin
		for i := 0; i < b.N; i++ {
			c <- message
		}
	}

	receiver := func() {
		defer wg.Done()
		<-begin
		for i := 0; i < b.N; i++ {
			<-c
		}
	}

	wg.Add(2)
	go sender()
	go receiver()
	b.StartTimer()
	close(begin) //tell the two go routines to begin
	wg.Wait()
}

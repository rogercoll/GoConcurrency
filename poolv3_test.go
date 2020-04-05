package main

import (
	"fmt"
	"net"
	"log"
	"sync"
	"time"
	"io/ioutil"
	"testing"
)

//WITHOUT POOL

//simulation of connection to service
func connectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

//to compare the performance we will start a new connection to the service for every request
func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()
		wg.Done()
		for {
			conn, err:= server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			connectToService()
			fmt.Fprintln(conn, "")
			conn.Close()
		}
	}()
	return &wg
}


//USING POOL

func warmServiceConnCache() *sync.Pool {
	p := &sync.Pool {
		New: connectToService,
	}
	for i := 0; i < 10; i ++ {
		p.Put(p.New())
	}
	return p
}

func startNetworkDaemonWithPool() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		connPool := warmServiceConnCache()
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()
		wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			svcConn := connPool.Get()
			fmt.Fprintln(conn, "")
			connPool.Put(svcConn)
			conn.Close()
		}
	}()
	return &wg
}

func init() {
	//daemonStarted := startNetworkDaemon()
	daemonStarted := startNetworkDaemonWithPool()
	daemonStarted.Wait()
}


//go test poolv3_test.go -benchtime=10s -bench=.
func BenchmarkNetworkRequestWithoutPool(b *testing.B) {
	for i:= 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}
		if _,err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v",err)
		}
		conn.Close()
	}
}

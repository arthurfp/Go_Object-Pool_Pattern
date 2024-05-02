package main

import (
	"flag"
	"fmt"
	"object-pool-go/pkg/pool"
	"time"
)

func main() {
	var maxOpenConns int
	var timeout int

	flag.IntVar(&maxOpenConns, "maxOpenConns", 5, "Maximum number of open connections in the pool")
	flag.IntVar(&timeout, "timeout", 30, "Timeout in seconds for waiting to borrow a connection")
	flag.Parse()

	fmt.Println("Object Pool Pattern in Go")

	connPool := pool.NewConnectionPool(maxOpenConns, time.Duration(timeout)*time.Second)

	// Demonstrate borrowing a connection
	conn, err := connPool.BorrowConnection()
	if err != nil {
		fmt.Println("Failed to borrow connection:", err)
		return
	}
	fmt.Println("Borrowed Connection ID:", conn.ID)

	// Release the connection
	connPool.ReleaseConnection(conn)
	fmt.Println("Connection released back to pool")
}

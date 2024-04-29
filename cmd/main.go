package main

import (
	"fmt"
	"object-pool-go/pkg/pool"
	"time"
)

func main() {
	fmt.Println("Object Pool Pattern in Go")

	connPool := pool.NewConnectionPool(2, 5*time.Second) // Pool size of 2 and timeout of 5 seconds

	// Borrow a connection
	conn, err := connPool.BorrowConnection()
	if err != nil {
		fmt.Println("Failed to borrow connection:", err)
		return
	}
	fmt.Println("Borrowed Connection ID:", conn.ID)

	// Simulate query
	if err := conn.Execute("SELECT * FROM dummy_table"); err != nil {
		fmt.Println("Failed to execute query:", err)
	}

	// Release the connection
	connPool.ReleaseConnection(conn)
	fmt.Println("Connection released back to pool")
}

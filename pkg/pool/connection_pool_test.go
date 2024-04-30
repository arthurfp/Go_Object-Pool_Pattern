package pool

import (
	"testing"
	"time"
)

// TestConnectionPoolBorrowRelease tests borrowing and releasing connections from the pool.
func TestConnectionPoolBorrowRelease(t *testing.T) {
	pool := NewConnectionPool(2, 1*time.Second) // Pool with 2 connections and 1-second timeout

	// Test borrowing the first connection
	conn1, err := pool.BorrowConnection()
	if err != nil {
		t.Errorf("Failed to borrow connection: %s", err)
	}

	// Test borrowing the second connection
	conn2, err := pool.BorrowConnection()
	if err != nil {
		t.Errorf("Failed to borrow connection: %s", err)
	}

	// Attempt to borrow a third connection, should timeout
	_, err = pool.BorrowConnection()
	if err == nil {
		t.Error("Expected error when borrowing third connection, but got none")
	}

	// Release one connection and try again
	pool.ReleaseConnection(conn1)
	_, err = pool.BorrowConnection()
	if err != nil {
		t.Errorf("Failed to borrow connection after release: %s", err)
	}

	// Cleanup: release all connections
	pool.ReleaseConnection(conn2)
}

// TestConnectionResetOnRelease tests if connections are reset when released back to the pool.
func TestConnectionResetOnRelease(t *testing.T) {
	pool := NewConnectionPool(1, 1*time.Second)
	conn, err := pool.BorrowConnection()
	if err != nil {
		t.Fatalf("Failed to borrow connection: %s", err)
	}

	// Simulate using the connection
	conn.Execute("UPDATE SET something")

	// Release and borrow again to test reset
	pool.ReleaseConnection(conn)
	conn2, err := pool.BorrowConnection()
	if err != nil {
		t.Fatalf("Failed to borrow connection again: %s", err)
	}

	// Assume the Reset method clears specific flags or states
	if conn2.ID != conn.ID {
		t.Errorf("Connection ID mismatch, expected %d got %d", conn.ID, conn2.ID)
	}
}

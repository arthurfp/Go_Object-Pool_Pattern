package pool

import (
	"object-pool-go/internal/database"
	"sync"
)

type ConnectionPool struct {
	mutex        sync.Mutex
	connections  []*database.Connection
	maxOpenConns int
}

func NewConnectionPool(maxOpenConns int) *ConnectionPool {
	return &ConnectionPool{
		connections:  make([]*database.Connection, 0, maxOpenConns),
		maxOpenConns: maxOpenConns,
	}
}

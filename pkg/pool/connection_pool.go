package pool

import (
	"errors"
	"object-pool-go/internal/database"
	"sync"
)

type ConnectionPool struct {
	mutex        sync.Mutex
	connections  []*database.Connection
	maxOpenConns int
}

func NewConnectionPool(maxOpenConns int) *ConnectionPool {
	pool := &ConnectionPool{
		connections:  make([]*database.Connection, 0, maxOpenConns),
		maxOpenConns: maxOpenConns,
	}
	// Pre-allocate connections
	for i := 0; i < maxOpenConns; i++ {
		pool.connections = append(pool.connections, &database.Connection{ID: i})
	}
	return pool
}

func (p *ConnectionPool) BorrowConnection() (*database.Connection, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.connections) > 0 {
		conn := p.connections[0]
		p.connections = p.connections[1:] // Remove the connection from the pool
		return conn, nil
	}
	return nil, errors.New("no connections available")
}

func (p *ConnectionPool) ReleaseConnection(conn *database.Connection) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Prepend the connection back to the pool
	p.connections = append([]*database.Connection{conn}, p.connections...)
}

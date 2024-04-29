package pool

import (
	"errors"
	"object-pool-go/internal/database"
	"sync"
	"time"
)

// ConnectionPool manages a pool of database connections.
type ConnectionPool struct {
	mutex        sync.Mutex
	connections  []*database.Connection
	maxOpenConns int
	timeout      time.Duration
}

// NewConnectionPool initializes a new connection pool with a specified size and timeout.
func NewConnectionPool(maxOpenConns int, timeout time.Duration) *ConnectionPool {
	pool := &ConnectionPool{
		connections:  make([]*database.Connection, 0, maxOpenConns),
		maxOpenConns: maxOpenConns,
		timeout:      timeout,
	}
	for i := 0; i < maxOpenConns; i++ {
		pool.connections = append(pool.connections, &database.Connection{ID: i})
	}
	return pool
}

// BorrowConnection attempts to borrow a connection from the pool, waiting if necessary until one is available or the timeout is exceeded.
func (p *ConnectionPool) BorrowConnection() (*database.Connection, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.connections) == 0 {
		select {
		case <-time.After(p.timeout):
			return nil, errors.New("timeout exceeded, no connections available")
		default:
		}
	}

	conn := p.connections[0]
	p.connections = p.connections[1:]
	return conn, nil
}

// ReleaseConnection returns a connection to the pool, resetting it for reuse.
func (p *ConnectionPool) ReleaseConnection(conn *database.Connection) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	conn.Reset()
	p.connections = append([]*database.Connection{conn}, p.connections...)
}

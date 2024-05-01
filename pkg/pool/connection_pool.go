package pool

import (
	"errors"
	"object-pool-go/internal/database"
	"sync"
	"time"
)

type ConnectionPool struct {
	mutex        sync.Mutex
	connections  []*database.Connection
	maxOpenConns int
	currentConns int
	timeout      time.Duration
}

func NewConnectionPool(maxOpenConns int, timeout time.Duration) *ConnectionPool {
	pool := &ConnectionPool{
		connections:  make([]*database.Connection, 0, maxOpenConns),
		maxOpenConns: maxOpenConns,
		currentConns: 0,
		timeout:      timeout,
	}
	// Pre-fill the pool with initial connections
	pool.expandPool(maxOpenConns / 2) // Start with half of the max capacity
	return pool
}

func (p *ConnectionPool) BorrowConnection() (*database.Connection, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Check for available connections or wait for timeout
	timeoutChan := time.After(p.timeout)
	for {
		if len(p.connections) > 0 {
			conn := p.connections[0]
			p.connections = p.connections[1:]
			if conn.HealthCheck() {
				return conn, nil
			}
		}

		// Expand the pool if under max capacity
		if p.currentConns < p.maxOpenConns {
			p.expandPool(1)
		}

		select {
		case <-timeoutChan:
			return nil, errors.New("timeout exceeded, no healthy connections available")
		default:
			time.Sleep(100 * time.Millisecond) // Prevent busy looping
		}
	}
}

func (p *ConnectionPool) ReleaseConnection(conn *database.Connection) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if conn.HealthCheck() {
		p.connections = append([]*database.Connection{conn}, p.connections...)
	} else {
		// Optionally replace an unhealthy connection
		if p.currentConns < p.maxOpenConns {
			p.expandPool(1)
		}
	}
}

func (p *ConnectionPool) expandPool(num int) {
	for i := 0; i < num; i++ {
		if p.currentConns >= p.maxOpenConns {
			return
		}
		p.connections = append(p.connections, &database.Connection{ID: p.currentConns})
		p.currentConns++
	}
}

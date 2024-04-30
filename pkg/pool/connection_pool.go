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
	timeout      time.Duration
}

func NewConnectionPool(maxOpenConns int, timeout time.Duration) *ConnectionPool {
	pool := &ConnectionPool{
		connections:  make([]*database.Connection, 0, maxOpenConns),
		maxOpenConns: maxOpenConns,
		timeout:      timeout,
	}
	for i := 0; i < maxOpenConns; i++ {
		conn := &database.Connection{ID: i}
		if conn.HealthCheck() {
			pool.connections = append(pool.connections, conn)
		}
	}
	return pool
}

func (p *ConnectionPool) BorrowConnection() (*database.Connection, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Wait until a connection becomes available or the timeout is reached
	timeoutChan := time.After(p.timeout)
	for {
		for len(p.connections) > 0 {
			conn := p.connections[0]
			p.connections = p.connections[1:]
			if conn.HealthCheck() {
				return conn, nil
			}
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
		// Replace unhealthy connection
		newConn := &database.Connection{ID: conn.ID}
		if newConn.HealthCheck() {
			p.connections = append([]*database.Connection{newConn}, p.connections...)
		}
	}
}

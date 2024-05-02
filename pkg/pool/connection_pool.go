package pool

import (
	"errors"
	"log"
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
	log.Println("Initializing new connection pool")
	pool.expandPool(maxOpenConns / 2) // Start with half of the max capacity
	return pool
}

func (p *ConnectionPool) BorrowConnection() (*database.Connection, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	timeoutChan := time.After(p.timeout)
	for {
		if len(p.connections) > 0 {
			conn := p.connections[0]
			p.connections = p.connections[1:]
			if conn.HealthCheck() {
				log.Printf("Borrowed connection ID %d", conn.ID)
				return conn, nil
			}
		}

		if p.currentConns < p.maxOpenConns {
			log.Println("Expanding pool due to high demand")
			p.expandPool(1)
		}

		select {
		case <-timeoutChan:
			log.Println("Failed to borrow connection: timeout exceeded")
			return nil, errors.New("timeout exceeded, no healthy connections available")
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (p *ConnectionPool) ReleaseConnection(conn *database.Connection) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if conn.HealthCheck() {
		p.connections = append([]*database.Connection{conn}, p.connections...)
		log.Printf("Returned connection ID %d to pool", conn.ID)
	} else {
		log.Printf("Connection ID %d failed health check, not returned to pool", conn.ID)
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
		log.Printf("Added new connection ID %d to pool", p.currentConns)
	}
}

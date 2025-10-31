package grpc

import (
	"errors"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
)

// Pool represents a connection pool for gRPC clients.
// It maintains multiple connections to improve throughput and reduce connection overhead.
type Pool struct {
	conns   []*grpc.ClientConn
	size    int
	counter uint64 // Atomic counter for round-robin
	mu      sync.RWMutex
	closed  bool
}

// PoolConfig holds configuration for creating a connection pool.
type PoolConfig struct {
	// Size is the number of connections in the pool
	Size int

	// Target is the gRPC server address (e.g., "localhost:50053")
	Target string

	// DialOptions are additional gRPC dial options
	DialOptions []grpc.DialOption
}

// NewPool creates a new gRPC connection pool.
// It establishes 'size' number of connections to the target server.
//
// Example:
//
//	pool, err := grpcpkg.NewPool(grpcpkg.PoolConfig{
//	    Size:   5,
//	    Target: "localhost:50053",
//	    DialOptions: []grpc.DialOption{
//	        grpc.WithTransportCredentials(insecure.NewCredentials()),
//	        grpc.WithUnaryInterceptor(loggingInterceptor),
//	    },
//	})
//	defer pool.Close()
//
//	conn := pool.Get()
//	client := pb.NewListingsServiceClient(conn)
func NewPool(config PoolConfig) (*Pool, error) {
	if config.Size <= 0 {
		return nil, errors.New("pool size must be greater than 0")
	}

	if config.Target == "" {
		return nil, errors.New("target address cannot be empty")
	}

	pool := &Pool{
		conns: make([]*grpc.ClientConn, config.Size),
		size:  config.Size,
	}

	// Establish connections
	for i := 0; i < config.Size; i++ {
		conn, err := grpc.Dial(config.Target, config.DialOptions...)
		if err != nil {
			// Close any connections that were successfully created
			for j := 0; j < i; j++ {
				pool.conns[j].Close()
			}
			return nil, err
		}
		pool.conns[i] = conn
	}

	return pool, nil
}

// Get returns a connection from the pool using round-robin strategy.
// This method is thread-safe.
func (p *Pool) Get() *grpc.ClientConn {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return nil
	}

	// Round-robin selection
	index := atomic.AddUint64(&p.counter, 1) % uint64(p.size)
	return p.conns[index]
}

// GetAll returns all connections in the pool.
// This is useful if you want to broadcast a request to all connections.
func (p *Pool) GetAll() []*grpc.ClientConn {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return nil
	}

	// Return a copy to avoid external modifications
	conns := make([]*grpc.ClientConn, p.size)
	copy(conns, p.conns)
	return conns
}

// Size returns the number of connections in the pool.
func (p *Pool) Size() int {
	return p.size
}

// Close closes all connections in the pool.
// After calling Close, Get() will return nil.
func (p *Pool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return errors.New("pool already closed")
	}

	p.closed = true

	var errs []error
	for i, conn := range p.conns {
		if err := conn.Close(); err != nil {
			errs = append(errs, err)
		}
		p.conns[i] = nil
	}

	if len(errs) > 0 {
		return errors.New("failed to close some connections")
	}

	return nil
}

// IsClosed returns true if the pool has been closed.
func (p *Pool) IsClosed() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.closed
}

// Reconnect attempts to reconnect a specific connection in the pool.
// This is useful if a connection becomes unhealthy.
//
// Example:
//
//	if err := pool.Reconnect(0, target, opts...); err != nil {
//	    log.Error("Failed to reconnect", err)
//	}
func (p *Pool) Reconnect(index int, target string, opts ...grpc.DialOption) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return errors.New("pool is closed")
	}

	if index < 0 || index >= p.size {
		return errors.New("index out of range")
	}

	// Close old connection
	if p.conns[index] != nil {
		p.conns[index].Close()
	}

	// Create new connection
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return err
	}

	p.conns[index] = conn
	return nil
}

// HealthCheck performs a basic health check on all connections.
// Returns the number of healthy connections.
//
// Note: This is a basic check using GetState(). For production,
// implement proper health checking using gRPC health checking protocol.
func (p *Pool) HealthCheck() int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.closed {
		return 0
	}

	healthy := 0
	for _, conn := range p.conns {
		if conn != nil {
			state := conn.GetState()
			// Consider Ready and Idle as healthy
			if state.String() == "READY" || state.String() == "IDLE" {
				healthy++
			}
		}
	}

	return healthy
}

// Stats returns statistics about the connection pool.
type PoolStats struct {
	Size           int
	Closed         bool
	HealthyConns   int
	RequestCounter uint64
}

// GetStats returns current statistics about the pool.
func (p *Pool) GetStats() PoolStats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return PoolStats{
		Size:           p.size,
		Closed:         p.closed,
		HealthyConns:   p.HealthCheck(),
		RequestCounter: atomic.LoadUint64(&p.counter),
	}
}

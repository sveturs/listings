package testing

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

// GRPCTestClient wraps a gRPC client connection with test-friendly utilities.
// It provides automatic timeout management, connection pooling, and cleanup.
type GRPCTestClient struct {
	client   pb.ListingsServiceClient
	conn     *grpc.ClientConn
	listener *bufconn.Listener
	server   *grpc.Server
	logger   zerolog.Logger
}

// GRPCTestClientConfig holds configuration for creating a test client
type GRPCTestClientConfig struct {
	// BufSize is the size of the in-memory buffer for bufconn (default: 1MB)
	BufSize int

	// DefaultTimeout is the default timeout for RPC calls (default: 5s)
	DefaultTimeout time.Duration

	// Logger is the logger to use (if nil, a noop logger is used)
	Logger *zerolog.Logger

	// DialOptions are additional gRPC dial options
	DialOptions []grpc.DialOption

	// ServerOptions are additional gRPC server options
	ServerOptions []grpc.ServerOption
}

// DefaultGRPCTestClientConfig returns a configuration with sensible defaults
func DefaultGRPCTestClientConfig() GRPCTestClientConfig {
	return GRPCTestClientConfig{
		BufSize:        1024 * 1024, // 1MB
		DefaultTimeout: 5 * time.Second,
		DialOptions:    []grpc.DialOption{},
		ServerOptions:  []grpc.ServerOption{},
	}
}

// NewGRPCTestClient creates a new test client connected to an in-memory gRPC server.
// This uses bufconn for efficient in-process testing without network overhead.
//
// The server parameter should be the gRPC server implementation to test.
//
// Example:
//
//	testClient := testing.NewGRPCTestClient(grpcServer, testing.DefaultGRPCTestClientConfig())
//	defer testClient.Close()
//
//	ctx := testClient.Context()
//	resp, err := testClient.Client().GetListing(ctx, &pb.GetListingRequest{ListingId: 1})
func NewGRPCTestClient(server pb.ListingsServiceServer, config GRPCTestClientConfig) (*GRPCTestClient, error) {
	// Apply defaults
	if config.BufSize <= 0 {
		config.BufSize = 1024 * 1024
	}
	if config.DefaultTimeout <= 0 {
		config.DefaultTimeout = 5 * time.Second
	}

	// Setup logger
	var logger zerolog.Logger
	if config.Logger != nil {
		logger = *config.Logger
	} else {
		logger = zerolog.Nop()
	}

	// Create in-memory listener
	listener := bufconn.Listen(config.BufSize)

	// Create gRPC server with options
	grpcServer := grpc.NewServer(config.ServerOptions...)
	pb.RegisterListingsServiceServer(grpcServer, server)

	// Start server in background
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error().Err(err).Msg("gRPC test server failed")
		}
	}()

	// Create client connection with bufconn dialer
	dialOpts := append(
		config.DialOptions,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
	)

	conn, err := grpc.DialContext(context.Background(), "bufnet", dialOpts...)
	if err != nil {
		grpcServer.Stop()
		listener.Close()
		return nil, fmt.Errorf("failed to dial test server: %w", err)
	}

	client := pb.NewListingsServiceClient(conn)

	return &GRPCTestClient{
		client:   client,
		conn:     conn,
		listener: listener,
		server:   grpcServer,
		logger:   logger,
	}, nil
}

// Client returns the underlying gRPC client for making RPC calls
func (c *GRPCTestClient) Client() pb.ListingsServiceClient {
	return c.client
}

// Context creates a new context with the default timeout.
// This is a convenience method for creating contexts for RPC calls.
func (c *GRPCTestClient) Context() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	return ctx
}

// ContextWithTimeout creates a new context with a custom timeout.
// Use this when you need a different timeout than the default.
func (c *GRPCTestClient) ContextWithTimeout(timeout time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx
}

// ContextWithCancel creates a new context with a cancel function.
// Remember to call the cancel function when done to avoid context leaks.
func (c *GRPCTestClient) ContextWithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

// Close gracefully shuts down the test client and server.
// This should be called in a defer statement after creating the client.
func (c *GRPCTestClient) Close() error {
	var errors []error

	// Close client connection
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close connection: %w", err))
		}
	}

	// Stop gRPC server gracefully
	if c.server != nil {
		c.server.GracefulStop()
	}

	// Close listener
	if c.listener != nil {
		if err := c.listener.Close(); err != nil {
			errors = append(errors, fmt.Errorf("failed to close listener: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors during cleanup: %v", errors)
	}

	return nil
}

// WaitForReady blocks until the connection is ready or the context is canceled.
// This is useful for ensuring the server is fully initialized before making calls.
func (c *GRPCTestClient) WaitForReady(ctx context.Context) error {
	// Check connection state
	state := c.conn.GetState()
	if state.String() == "READY" {
		return nil
	}

	// Wait for state to become ready
	c.conn.WaitForStateChange(ctx, state)

	newState := c.conn.GetState()
	if newState.String() != "READY" {
		return fmt.Errorf("connection not ready: state=%s", newState)
	}

	return nil
}

// IsHealthy returns true if the connection is in a healthy state (READY or IDLE)
func (c *GRPCTestClient) IsHealthy() bool {
	state := c.conn.GetState()
	return state.String() == "READY" || state.String() == "IDLE"
}

// GetConnectionState returns the current state of the gRPC connection
func (c *GRPCTestClient) GetConnectionState() string {
	return c.conn.GetState().String()
}

// GRPCTestClientPool manages multiple test clients for parallel testing.
// This is useful for load testing or simulating multiple concurrent clients.
type GRPCTestClientPool struct {
	clients []*GRPCTestClient
	size    int
	logger  zerolog.Logger
}

// NewGRPCTestClientPool creates a pool of test clients for parallel testing.
//
// Example:
//
//	pool := testing.NewGRPCTestClientPool(5, grpcServer, config)
//	defer pool.CloseAll()
//
//	// Run tests in parallel
//	t.Run("Parallel", func(t *testing.T) {
//	    for i := 0; i < pool.Size(); i++ {
//	        client := pool.Get(i)
//	        t.Run(fmt.Sprintf("Client%d", i), func(t *testing.T) {
//	            t.Parallel()
//	            // Use client for testing
//	        })
//	    }
//	})
func NewGRPCTestClientPool(size int, server pb.ListingsServiceServer, config GRPCTestClientConfig) (*GRPCTestClientPool, error) {
	if size <= 0 {
		return nil, fmt.Errorf("pool size must be positive, got %d", size)
	}

	var logger zerolog.Logger
	if config.Logger != nil {
		logger = *config.Logger
	} else {
		logger = zerolog.Nop()
	}

	clients := make([]*GRPCTestClient, size)
	for i := 0; i < size; i++ {
		client, err := NewGRPCTestClient(server, config)
		if err != nil {
			// Cleanup already created clients
			for j := 0; j < i; j++ {
				clients[j].Close()
			}
			return nil, fmt.Errorf("failed to create client %d: %w", i, err)
		}
		clients[i] = client
	}

	return &GRPCTestClientPool{
		clients: clients,
		size:    size,
		logger:  logger,
	}, nil
}

// Get returns the client at the specified index
func (p *GRPCTestClientPool) Get(index int) *GRPCTestClient {
	if index < 0 || index >= p.size {
		p.logger.Warn().Int("index", index).Int("size", p.size).Msg("index out of range")
		return nil
	}
	return p.clients[index]
}

// GetAll returns all clients in the pool
func (p *GRPCTestClientPool) GetAll() []*GRPCTestClient {
	return p.clients
}

// Size returns the number of clients in the pool
func (p *GRPCTestClientPool) Size() int {
	return p.size
}

// CloseAll closes all clients in the pool
func (p *GRPCTestClientPool) CloseAll() error {
	var errors []error
	for i, client := range p.clients {
		if err := client.Close(); err != nil {
			errors = append(errors, fmt.Errorf("client %d: %w", i, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to close some clients: %v", errors)
	}

	return nil
}

// AreAllHealthy returns true if all clients in the pool are healthy
func (p *GRPCTestClientPool) AreAllHealthy() bool {
	for _, client := range p.clients {
		if !client.IsHealthy() {
			return false
		}
	}
	return true
}

// CountHealthy returns the number of healthy clients in the pool
func (p *GRPCTestClientPool) CountHealthy() int {
	count := 0
	for _, client := range p.clients {
		if client.IsHealthy() {
			count++
		}
	}
	return count
}

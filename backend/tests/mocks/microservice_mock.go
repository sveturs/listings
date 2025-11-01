// Package main provides a mock gRPC microservice for testing
// backend/tests/mocks/microservice_mock.go
//
// Run with: go run tests/mocks/microservice_mock.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	grpcPort    = ":50051"
	controlPort = ":50052" // HTTP API for test control
)

// MockMode defines behavior mode
type MockMode string

const (
	ModeNormal  MockMode = "normal"  // Normal responses
	ModeSlow    MockMode = "slow"    // Slow responses (configurable delay)
	ModeError   MockMode = "error"   // Always error
	ModePartial MockMode = "partial" // Partial failures (configurable rate)
)

// Config holds mock service configuration
type Config struct {
	Mode        MockMode      `json:"mode"`
	Delay       time.Duration `json:"delay"`        // For slow mode
	FailureRate int           `json:"failure_rate"` // For partial mode (0-100)
	mu          sync.RWMutex
}

var config = &Config{
	Mode:        ModeNormal,
	Delay:       1 * time.Second,
	FailureRate: 50,
}

// ListingsServiceServer implements the gRPC service
type ListingsServiceServer struct {
	pb.UnimplementedListingsServiceServer
}

// GetListing returns a listing by ID
func (s *ListingsServiceServer) GetListing(ctx context.Context, req *pb.GetListingRequest) (*pb.GetListingResponse, error) {
	// Check for special test IDs
	switch req.Id {
	case 999:
		// Timeout test: delay 1 second
		time.Sleep(1 * time.Second)
		return nil, status.Error(codes.DeadlineExceeded, "request timeout")

	case 777:
		// Error test: return internal error
		return nil, status.Error(codes.Internal, "internal server error")

	case 888:
		// Unavailable test
		return nil, status.Error(codes.Unavailable, "service unavailable")
	}

	// Apply current mode behavior
	config.mu.RLock()
	mode := config.Mode
	delay := config.Delay
	failureRate := config.FailureRate
	config.mu.RUnlock()

	switch mode {
	case ModeSlow:
		// Delay response
		select {
		case <-ctx.Done():
			return nil, status.Error(codes.Canceled, "request canceled")
		case <-time.After(delay):
			// Continue to normal response
		}

	case ModeError:
		// Always return error
		return nil, status.Error(codes.Internal, "mock service configured to fail")

	case ModePartial:
		// Random failures based on failure rate
		if int(req.Id)%100 < failureRate {
			return nil, status.Error(codes.Internal, "partial failure mode")
		}

	case ModeNormal:
		// Normal operation - continue
	}

	// Normal response
	return &pb.GetListingResponse{
		Listing: &pb.Listing{
			Id:          req.Id,
			Title:       fmt.Sprintf("Test Listing %d", req.Id),
			Description: "This is a mock listing for testing",
			Price:       1000,
			Currency:    "USD",
			Status:      "active",
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
		},
	}, nil
}

// CreateListing creates a new listing
func (s *ListingsServiceServer) CreateListing(ctx context.Context, req *pb.CreateListingRequest) (*pb.CreateListingResponse, error) {
	// Apply mode behavior
	config.mu.RLock()
	mode := config.Mode
	config.mu.RUnlock()

	if mode == ModeError {
		return nil, status.Error(codes.Internal, "mock service configured to fail")
	}

	return &pb.CreateListingResponse{
		Listing: &pb.Listing{
			Id:          123,
			Title:       req.Title,
			Description: req.Description,
			Price:       req.Price,
			Currency:    req.Currency,
			Status:      "active",
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
		},
	}, nil
}

// UpdateListing updates a listing
func (s *ListingsServiceServer) UpdateListing(ctx context.Context, req *pb.UpdateListingRequest) (*pb.UpdateListingResponse, error) {
	config.mu.RLock()
	mode := config.Mode
	config.mu.RUnlock()

	if mode == ModeError {
		return nil, status.Error(codes.Internal, "mock service configured to fail")
	}

	return &pb.UpdateListingResponse{
		Listing: &pb.Listing{
			Id:          req.Id,
			Title:       req.Title,
			Description: req.Description,
			Price:       req.Price,
			Currency:    req.Currency,
			Status:      "active",
			UpdatedAt:   time.Now().Unix(),
		},
	}, nil
}

// DeleteListing soft-deletes a listing
func (s *ListingsServiceServer) DeleteListing(ctx context.Context, req *pb.DeleteListingRequest) (*pb.DeleteListingResponse, error) {
	config.mu.RLock()
	mode := config.Mode
	config.mu.RUnlock()

	if mode == ModeError {
		return nil, status.Error(codes.Internal, "mock service configured to fail")
	}

	return &pb.DeleteListingResponse{
		Success: true,
	}, nil
}

// SearchListings searches for listings
func (s *ListingsServiceServer) SearchListings(ctx context.Context, req *pb.SearchListingsRequest) (*pb.SearchListingsResponse, error) {
	config.mu.RLock()
	mode := config.Mode
	delay := config.Delay
	config.mu.RUnlock()

	if mode == ModeSlow {
		time.Sleep(delay)
	}

	if mode == ModeError {
		return nil, status.Error(codes.Internal, "mock service configured to fail")
	}

	// Return mock results
	return &pb.SearchListingsResponse{
		Listings: []*pb.Listing{
			{
				Id:          1,
				Title:       "Search Result 1",
				Description: "Mock search result",
				Price:       1000,
				Currency:    "USD",
				Status:      "active",
			},
			{
				Id:          2,
				Title:       "Search Result 2",
				Description: "Mock search result",
				Price:       2000,
				Currency:    "USD",
				Status:      "active",
			},
		},
		TotalCount: 2,
	}, nil
}

// ListListings lists listings with pagination
func (s *ListingsServiceServer) ListListings(ctx context.Context, req *pb.ListListingsRequest) (*pb.ListListingsResponse, error) {
	config.mu.RLock()
	mode := config.Mode
	config.mu.RUnlock()

	if mode == ModeError {
		return nil, status.Error(codes.Internal, "mock service configured to fail")
	}

	return &pb.ListListingsResponse{
		Listings: []*pb.Listing{
			{
				Id:          1,
				Title:       "Listing 1",
				Description: "Mock listing",
				Price:       1000,
				Currency:    "USD",
				Status:      "active",
			},
		},
		TotalCount: 1,
	}, nil
}

// HTTP Control API handlers

func handleConfigUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newConfig struct {
		Mode        string `json:"mode"`
		Delay       string `json:"delay"`
		FailureRate int    `json:"failure_rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	config.mu.Lock()
	defer config.mu.Unlock()

	// Update mode
	if newConfig.Mode != "" {
		config.Mode = MockMode(newConfig.Mode)
	}

	// Update delay
	if newConfig.Delay != "" {
		delay, err := time.ParseDuration(newConfig.Delay)
		if err == nil {
			config.Delay = delay
		}
	}

	// Update failure rate
	if newConfig.FailureRate > 0 && newConfig.FailureRate <= 100 {
		config.FailureRate = newConfig.FailureRate
	}

	log.Printf("Config updated: mode=%s, delay=%v, failure_rate=%d%%",
		config.Mode, config.Delay, config.FailureRate)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	config.mu.RLock()
	defer config.mu.RUnlock()

	response := map[string]interface{}{
		"mode":         string(config.Mode),
		"delay":        config.Delay.String(),
		"failure_rate": config.FailureRate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func main() {
	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterListingsServiceServer(grpcServer, &ListingsServiceServer{})

		log.Printf("Mock gRPC server listening on %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Start HTTP control server
	http.HandleFunc("/control/config", handleConfigUpdate)
	http.HandleFunc("/control/status", handleGetConfig)
	http.HandleFunc("/health", handleHealth)

	log.Printf("Mock control HTTP server listening on %s", controlPort)
	log.Printf("Available modes: normal, slow, error, partial")
	log.Printf("Control endpoint: POST http://localhost%s/control/config", controlPort)
	log.Printf("\nExample: curl -X POST http://localhost%s/control/config -d '{\"mode\":\"slow\",\"delay\":\"1s\"}'", controlPort)

	if err := http.ListenAndServe(controlPort, nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

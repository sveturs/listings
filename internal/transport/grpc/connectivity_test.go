package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPCConnectivity(t *testing.T) {
	t.Skip("Skipping gRPC connectivity test - requires fully migrated database with fixtures")

	if testing.Short() {
		t.Skip("Skipping gRPC connectivity test in short mode - requires running server")
	}

	// Connect to the gRPC server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "Failed to connect to gRPC server")
	defer conn.Close()

	client := pb.NewListingsServiceClient(conn)

	// Test ListListings (should work even with no data)
	resp, err := client.ListListings(ctx, &pb.ListListingsRequest{
		Limit:  10,
		Offset: 0,
	})

	require.NoError(t, err, "ListListings should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	require.GreaterOrEqual(t, resp.Total, int32(0), "Total should be >= 0")

	t.Logf("gRPC connectivity test passed: got %d listings", resp.Total)
}

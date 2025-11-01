package grpc

import (
	"context"
	"testing"
	"time"

	listingsv1 "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPCConnectivity(t *testing.T) {
	// Connect to the gRPC server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "localhost:50053",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err, "Failed to connect to gRPC server")
	defer conn.Close()

	client := listingsv1.NewListingsServiceClient(conn)

	// Test ListListings (should work even with no data)
	resp, err := client.ListListings(ctx, &listingsv1.ListListingsRequest{
		Limit:  10,
		Offset: 0,
	})

	require.NoError(t, err, "ListListings should succeed")
	require.NotNil(t, resp, "Response should not be nil")
	require.GreaterOrEqual(t, resp.Total, int32(0), "Total should be >= 0")

	t.Logf("gRPC connectivity test passed: got %d listings", resp.Total)
}

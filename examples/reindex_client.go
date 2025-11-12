package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

func main() {
	// Get gRPC server address from environment or use default
	grpcAddr := os.Getenv("GRPC_HOST")
	if grpcAddr == "" {
		grpcAddr = "localhost:50051"
	}

	fmt.Printf("🔌 Connecting to gRPC server at %s\n", grpcAddr)

	// Connect to gRPC server
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewListingsServiceClient(conn)

	// Create context with timeout (reindexing can take a while)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Parse command line arguments
	sourceType := ""
	batchSize := int32(1000)

	if len(os.Args) > 1 {
		sourceType = os.Args[1]
	}
	if len(os.Args) > 2 {
		fmt.Sscanf(os.Args[2], "%d", &batchSize)
	}

	// Build request
	req := &pb.ReindexAllRequest{}
	if sourceType != "" {
		req.SourceType = &sourceType
		fmt.Printf("📦 Source Type: %s\n", sourceType)
	} else {
		fmt.Println("📦 Source Type: ALL (b2c + c2c)")
	}
	if batchSize != 1000 {
		req.BatchSize = &batchSize
	}
	fmt.Printf("📦 Batch Size: %d\n", batchSize)
	fmt.Println()

	// Call ReindexAll
	fmt.Println("🚀 Starting reindexing...")
	startTime := time.Now()

	resp, err := client.ReindexAll(ctx, req)
	if err != nil {
		log.Fatalf("❌ ReindexAll failed: %v", err)
	}

	elapsed := time.Since(startTime)

	// Print results
	fmt.Println()
	fmt.Println("✅ Reindexing completed successfully!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("📊 Total Indexed:    %d products\n", resp.TotalIndexed)
	fmt.Printf("❌ Total Failed:     %d products\n", resp.TotalFailed)
	fmt.Printf("⏱️  Duration:         %d seconds (%.2f minutes)\n", resp.DurationSeconds, float64(resp.DurationSeconds)/60)
	fmt.Printf("🕐 Client Elapsed:   %.2f seconds\n", elapsed.Seconds())
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if resp.TotalFailed > 0 && len(resp.Errors) > 0 {
		fmt.Println()
		fmt.Println("⚠️  Sample Errors (first 10):")
		for i, errMsg := range resp.Errors {
			fmt.Printf("  %d. %s\n", i+1, errMsg)
		}
	}

	// Calculate success rate
	total := resp.TotalIndexed + resp.TotalFailed
	if total > 0 {
		successRate := float64(resp.TotalIndexed) / float64(total) * 100
		fmt.Println()
		fmt.Printf("📈 Success Rate: %.2f%%\n", successRate)
	}

	fmt.Println()
	fmt.Println("✨ Done!")
}

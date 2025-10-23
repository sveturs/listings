package grpcclient

import (
	"context"

	pb "backend/pkg/grpc/delivery/v1"
)

// Adapter адаптирует gRPC клиент к интерфейсу сервиса
// Это позволяет service слою использовать gRPC клиент через простой интерфейс
type Adapter struct {
	client *Client
}

// NewAdapter создает новый адаптер
func NewAdapter(client *Client) *Adapter {
	return &Adapter{
		client: client,
	}
}

// CreateShipmentGRPC создает отправление через gRPC
func (a *Adapter) CreateShipmentGRPC(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error) {
	return a.client.CreateShipment(ctx, req)
}

// GetShipmentGRPC получает отправление через gRPC
func (a *Adapter) GetShipmentGRPC(ctx context.Context, req *pb.GetShipmentRequest) (*pb.GetShipmentResponse, error) {
	return a.client.GetShipment(ctx, req)
}

// TrackShipmentGRPC отслеживает отправление через gRPC
func (a *Adapter) TrackShipmentGRPC(ctx context.Context, req *pb.TrackShipmentRequest) (*pb.TrackShipmentResponse, error) {
	return a.client.TrackShipment(ctx, req)
}

// CancelShipmentGRPC отменяет отправление через gRPC
func (a *Adapter) CancelShipmentGRPC(ctx context.Context, req *pb.CancelShipmentRequest) (*pb.CancelShipmentResponse, error) {
	return a.client.CancelShipment(ctx, req)
}

// CalculateRateGRPC рассчитывает стоимость через gRPC
func (a *Adapter) CalculateRateGRPC(ctx context.Context, req *pb.CalculateRateRequest) (*pb.CalculateRateResponse, error) {
	return a.client.CalculateRate(ctx, req)
}

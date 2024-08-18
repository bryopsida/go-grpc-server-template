package increment

import (
	"context"
	"log/slog"

	api_v1 "github.com/bryopsida/go-grpc-server-template/api/v1"
	"github.com/bryopsida/go-grpc-server-template/interfaces"
)

// IncrementServiceImpl is the implementation of IncrementServiceServer
type IncrementServiceImpl struct {
	api_v1.UnimplementedIncrementServiceServer
	repo   interfaces.INumberRepository
	bucket string
}

// NewIncrementService creates a new IncrementServiceImpl
// - repo: INumberRepository number repository
// - bucket: string bucket name
func NewIncrementService(repo interfaces.INumberRepository, bucket string) *IncrementServiceImpl {
	return &IncrementServiceImpl{
		repo:   repo,
		bucket: bucket,
	}
}

// Increment increments the number in the bucket
// - ctx: context.Context context
// - req: *api_v1.IncrementRequest request
// Returns *api_v1.IncrementResponse response
func (s *IncrementServiceImpl) Increment(ctx context.Context, req *api_v1.IncrementRequest) (*api_v1.IncrementResponse, error) {
	number, err := s.repo.FindByID(s.bucket)
	if err != nil {
		slog.Info("Bucket not found, creating new bucket", "bucket", s.bucket)
		number = &interfaces.Number{ID: s.bucket, Number: 0}
	}
	slog.Info("Incrementing number", "number", number.Number)
	number.Number++
	slog.Info("Saving number", "number", number.Number)
	saveErr := s.repo.Save(*number)
	if saveErr != nil {
		slog.Error("Error saving number", "error", saveErr)
		return nil, saveErr
	}

	resp := &api_v1.IncrementResponse{Value: number.Number}
	slog.Info("Returning incremented number", "number", resp.Value)
	return resp, nil
}

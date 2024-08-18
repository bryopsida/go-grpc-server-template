package increment

import (
	"context"

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
		number = &interfaces.Number{ID: s.bucket, Number: 0}
	}
	number.Number++

	err = s.repo.Save(*number)
	if err != nil {
		return nil, err
	}

	return &api_v1.IncrementResponse{Value: number.Number}, nil
}

package mngt

import (
	domain "github.com/kenriortega/ngonx/internal/mngt/domain"
)

type IMngtService interface {
	ListEnpoints() ([]domain.Endpoint, error)
	RegisterEnpoint(domain.Endpoint) error
}

type MngtService struct {
	repo domain.IEndpoint
}

// NewMngtService return new MngtService
func NewMngtService(repository domain.IEndpoint) MngtService {
	return MngtService{repo: repository}
}

func (s MngtService) ListEnpoints() ([]domain.Endpoint, error) {
	endpoints, err := s.repo.ListEnpoints()
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (s MngtService) RegisterEnpoint(endpoint domain.Endpoint) error {

	err := s.repo.RegisterEnpoint(endpoint)
	if err != nil {
		return err
	}

	return nil
}

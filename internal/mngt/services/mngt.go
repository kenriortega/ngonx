package mngt

import (
	domain "github.com/kenriortega/ngonx/internal/mngt/domain"
)

type IMngtService interface {
	ListEndpoints() ([]domain.Endpoint, error)
	RegisterEndpoint(domain.Endpoint) error
	UpdateEndpoint(domain.Endpoint) error
}

type MngtService struct {
	repo domain.IEndpoint
}

// NewMngtService return new MngtService
func NewMngtService(repository domain.IEndpoint) MngtService {
	return MngtService{repo: repository}
}

func (s MngtService) ListEndpoints() ([]domain.Endpoint, error) {
	endpoints, err := s.repo.ListEndpoints()

	if err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (s MngtService) RegisterEndpoint(endpoint domain.Endpoint) error {

	err := s.repo.RegisterEndpoint(endpoint)
	if err != nil {
		return err
	}

	return nil
}

func (s MngtService) UpdateEndpoint(endpoint domain.Endpoint) error {

	err := s.repo.UpdateEndpoint(endpoint)
	if err != nil {
		return err
	}

	return nil
}

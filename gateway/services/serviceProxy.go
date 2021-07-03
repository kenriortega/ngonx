package gateway

import (
	domain "egosystem.org/micros/gateway/domain"
)

type ProxyService interface {
	SaveSecretKEY(string, string, string) error
}

type DefaultProxyService struct {
	repo domain.ProxyRepository
}

func NewProxyService(repository domain.ProxyRepository) DefaultProxyService {
	return DefaultProxyService{repo: repository}
}

func (s DefaultProxyService) SaveSecretKEY(engine, key, apikey string) (string, error) {

	err := s.repo.SaveKEY(engine, key, apikey)
	if err != nil {
		return "failed", err
	}
	return "ok", nil
}
